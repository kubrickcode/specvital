---
title: Data Retention-Based Cleanup Service
description: ADR for implementing automated data retention cleanup as a standalone Railway Cron binary
---

# ADR-11: Data Retention-Based Cleanup Service

> 🇰🇷 [한국어 버전](/ko/adr/worker/11-data-retention-cleanup.md)

| Date       | Author       | Repos         |
| ---------- | ------------ | ------------- |
| 2026-02-02 | @KubrickCode | worker, infra |

## Context

### Data Growth Problem

The Specvital platform accumulates temporal data across multiple tables:

| Data Type          | Table                   | Growth Pattern          | Retention Need       |
| ------------------ | ----------------------- | ----------------------- | -------------------- |
| Analysis History   | `user_analysis_history` | Per-commit analysis     | Tier-based retention |
| Spec Documents     | `spec_documents`        | Per-SpecView generation | Tier-based retention |
| Analyses           | `analyses`              | Shared analysis records | Orphan cleanup       |
| Usage Events       | `usage_events`          | Per-operation (ADR-13)  | Audit + cleanup      |
| Quota Reservations | `quota_reservations`    | Per-concurrent request  | Short-lived (TTL)    |

Without active cleanup, storage costs grow unbounded and query performance degrades.

### Retention Policy Requirements

Subscription plans define tier-based retention periods:

| Tier       | Retention Period | Rationale               |
| ---------- | ---------------- | ----------------------- |
| Free       | 30 days          | Cost control            |
| Pro        | 90 days          | Standard business needs |
| Pro Plus   | 180 days         | Extended history        |
| Enterprise | Unlimited (NULL) | Compliance requirements |

### Downgrade Fairness Problem

**Challenge**: When users downgrade from Pro to Free, should their Pro-tier data (90-day retention) be deleted immediately?

**Answer**: No. Data created under a higher-tier plan should retain its original retention period. This requires capturing retention policy at record creation time, not at cleanup time.

### Architectural Context

The Scheduler service removal ([ADR-22](/en/adr/22-scheduler-removal-railway-cron.md)) eliminated the centralized cron job runner. Cleanup tasks previously embedded in the Scheduler now require Railway Cron deployment.

## Decision

**Implement data retention cleanup as an independent binary (`cmd/retention-cleanup`) triggered by Railway Cron at 3:00 AM UTC daily.**

### Creation-Time Retention Snapshot

Store `retention_days_at_creation` when records are created:

```sql
ALTER TABLE user_analysis_history ADD COLUMN retention_days_at_creation integer;
ALTER TABLE spec_documents ADD COLUMN retention_days_at_creation integer;

-- Constraint: positive or NULL
CHECK ((retention_days_at_creation IS NULL) OR (retention_days_at_creation > 0))

-- Partial index for cleanup queries
CREATE INDEX idx_uah_retention_cleanup ON user_analysis_history (created_at)
WHERE (retention_days_at_creation IS NOT NULL);
```

**Value semantics**:

| Value        | Meaning                                 |
| ------------ | --------------------------------------- |
| NULL         | Unlimited retention (enterprise/legacy) |
| Positive int | Days until eligible for deletion        |

### Two-Phase Cleanup Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                    Cleanup Execution Order                       │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1a: Delete expired user_analysis_history                  │
│            WHERE created_at + retention_days < NOW()             │
│                                                                  │
│  Phase 1b: Delete expired spec_documents                         │
│            WHERE created_at + retention_days < NOW()             │
│                                                                  │
│  Phase 2:  Delete orphaned analyses                              │
│            WHERE no remaining user_analysis_history refs         │
│            AND created_at < NOW() - 1 day (grace period)         │
└─────────────────────────────────────────────────────────────────┘
```

### Binary Architecture

```
cmd/retention-cleanup/
├── main.go                    # Entry point, run-to-completion
└── (bootstrap via internal/)

src/internal/
├── domain/retention/
│   ├── errors.go              # Domain errors
│   ├── policy.go              # Expiration calculation
│   └── repository.go          # CleanupRepository interface
├── usecase/retention/
│   └── cleanup.go             # CleanupUseCase orchestration
└── adapter/repository/postgres/
    └── retention.go           # PostgreSQL implementation
```

### CleanupRepository Interface

```go
type CleanupRepository interface {
    DeleteExpiredUserAnalysisHistory(ctx context.Context, batchSize int) (DeleteResult, error)
    DeleteExpiredSpecDocuments(ctx context.Context, batchSize int) (DeleteResult, error)
    DeleteOrphanedAnalyses(ctx context.Context, batchSize int) (DeleteResult, error)
}

type DeleteResult struct {
    DeletedCount int64
    HasMore      bool
}
```

### Railway Configuration

```json
{
  "$schema": "https://railway.com/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "infra/retention-cleanup/Dockerfile"
  },
  "deploy": {
    "cronSchedule": "0 3 * * *",
    "restartPolicyType": "NEVER"
  }
}
```

## Options Considered

### Option A: Individual Binary + Railway Cron (Selected)

Standalone Go binary executed daily via Railway Cron.

**Pros**:

- Failure isolation from other workloads
- Independent deployment and configuration
- No 24/7 running cost (pay per execution)
- Railway handles scheduling reliability

**Cons**:

- Cold start latency (~10-30 seconds)
- Railway platform dependency
- Another binary to maintain

### Option B: Embedded in Scheduler Service (Previous)

Cleanup logic as cron job within centralized Scheduler.

**Pros**:

- Single service to monitor
- Warm execution (no cold start)

**Cons**:

- Scheduler removed per [ADR-22](/en/adr/22-scheduler-removal-railway-cron.md)
- 24/7 running cost for daily job

### Option C: PostgreSQL pg_cron Extension

Database-level scheduled job.

**Pros**:

- No application binary needed
- Native PostgreSQL solution

**Cons**:

- pg_cron not available on Railway Postgres
- Complex tier-aware logic in SQL
- No application-level logging

### Option D: Current-Plan Lookup at Cleanup Time

Query user's current plan during cleanup instead of storing at creation.

**Pros**:

- No schema changes needed
- Always reflects current tier

**Cons**:

- **Unfair on downgrade**: Pro data deleted immediately when user downgrades to Free
- User loses data they paid to retain

## Consequences

### Positive

| Area                 | Benefit                                  |
| -------------------- | ---------------------------------------- |
| Storage Optimization | Prevents unbounded table growth          |
| Query Performance    | Smaller tables maintain index efficiency |
| Cost Control         | Per-execution billing (no 24/7 process)  |
| Downgrade Fairness   | Existing data retains original retention |
| Compliance Readiness | Automated data lifecycle management      |
| Audit Trail          | Deletion counts logged per execution     |

### Negative

| Area                   | Trade-off                                   |
| ---------------------- | ------------------------------------------- |
| Scheduling Granularity | Daily minimum (Railway Cron limitation)     |
| Cold Start             | 10-30 second startup overhead per execution |
| Platform Dependency    | Scheduling tied to Railway                  |
| Schema Addition        | New column on high-write tables             |

### Technical Notes

- **Deletion order**: Phase 1 before Phase 2 (foreign key awareness)
- **Batch processing**: Configurable batch size with sleep intervals
- **Timeout**: Railway execution timeout set to 30 minutes
- **Monitoring**: Total records deleted per table logged for trending

## Configuration

```
DATABASE_URL=postgres://...
RETENTION_TIMEOUT=30m
RETENTION_BATCH_SIZE=1000
RETENTION_BATCH_SLEEP=100ms
```

## References

- [ADR-22: Scheduler Removal and Railway Cron Migration](/en/adr/22-scheduler-removal-railway-cron.md) - Parent architectural decision
- [Worker ADR-08: SpecView Worker Binary Separation](/en/adr/worker/08-specview-worker-separation.md) - Binary separation pattern
- [ADR-13: Billing and Quota Architecture](/en/adr/13-billing-quota-architecture.md) - Defines `retention_days` in plans
- Commits: `7eb93aa`, `878d87c`, `5e8c05a`, `792f106`, `6e03a7f` (2026-02-02)
