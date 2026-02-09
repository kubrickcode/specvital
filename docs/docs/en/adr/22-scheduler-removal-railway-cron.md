---
title: Scheduler Removal and Railway Cron Migration
description: ADR on removing the Scheduler service and migrating to Railway Cron + individual binaries
---

# ADR-22: Scheduler Removal and Railway Cron Migration

> 🇰🇷 [한국어 버전](/ko/adr/22-scheduler-removal-railway-cron.md)

| Date       | Author       | Repos         |
| ---------- | ------------ | ------------- |
| 2026-02-02 | @KubrickCode | worker, infra |

## Status

**Accepted** - Supersedes [ADR-01: Scheduled Re-collection](/en/adr/worker/01-scheduled-recollection.md)

## Context

### Original Architecture Problem

The Scheduler service ([ADR-01](/en/adr/worker/01-scheduled-recollection.md)) was designed to pre-analyze repositories for instant user responses. However, operational data revealed fundamental cost-benefit issues:

| Metric            | Expected     | Actual                      |
| ----------------- | ------------ | --------------------------- |
| Analysis time     | 30+ seconds  | ~5 seconds                  |
| Pre-compute value | High         | Low (5s wait is acceptable) |
| Data freshness    | Maintainable | Impossible for active repos |
| Storage growth    | Controlled   | Rapid (unviewed results)    |
| 24/7 running cost | Justified    | Excessive for utility       |

### Why Pre-Compute Failed

1. **Low value proposition**: 5-second analysis time is acceptable for users
2. **Freshness impossible**: Active repositories have frequent commits, making pre-computed results immediately stale
3. **Database bloat**: Unviewed analysis results accumulated rapidly
4. **Cost inefficiency**: 24/7 scheduler running cost exceeded actual utility

### Scheduler Architecture Overhead

The Scheduler service introduced significant complexity:

- **Distributed locking**: PostgreSQL-based lock for single-instance guarantee
- **go-cron internal scheduling**: In-process cron job management
- **24/7 running cost**: Always-on service with minimal actual work
- **Failure coupling**: Scheduler failure affected all scheduled jobs

## Decision

**Remove the Scheduler service entirely. Migrate to Railway Cron triggers with individual single-purpose binaries.**

### New Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Before (Scheduler Service)                │
├─────────────────────────────────────────────────────────────┤
│  cmd/scheduler/    → 24/7 running, internal go-cron         │
│                      ├── Auto-refresh cron job              │
│                      ├── Distributed lock (PostgreSQL)      │
│                      └── Cleanup jobs (embedded)            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    After (Railway Cron)                      │
├─────────────────────────────────────────────────────────────┤
│  cmd/analyzer/         → Queue consumer (River, ON_FAILURE) │
│  cmd/spec-generator/   → Queue consumer (River, ON_FAILURE) │
│  cmd/retention-cleanup/→ Cron binary (Railway, "0 3 * * *") │
│  cmd/enqueue/          → Manual utility                     │
└─────────────────────────────────────────────────────────────┘
```

### Railway Cron Configuration

Each cron job is a separate Railway service with its own configuration:

**retention-cleanup/railway.json**:

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

**Key configuration points**:

- `cronSchedule`: Standard cron expression for scheduling
- `restartPolicyType: NEVER`: Binary runs to completion and exits
- Separate Dockerfile per service for clean build isolation

## Options Considered

### Option A: Railway Cron + Individual Binaries (Selected)

**How It Works:**

1. Each periodic task becomes an independent Go binary
2. Railway triggers binary execution via cron expression
3. Binary runs to completion and exits (no 24/7 process)
4. No distributed locking needed (Railway handles single execution)

**Pros:**

- No 24/7 running cost for cron scheduling
- No distributed lock complexity (platform manages this)
- Per-job cost visibility in Railway dashboard
- Railway handles scheduling reliability and retries
- Simpler deployment (just a binary that exits)
- Independent scaling and configuration per job

**Cons:**

- Cold start latency for each execution
- Railway platform dependency for scheduling
- Requires Railway IaC for each cron job

### Option B: Keep Scheduler with Reduced Scope

**Description:**

Maintain Scheduler but remove auto-refresh; keep only cleanup jobs.

**Pros:**

- Minimal code changes
- Existing monitoring and alerting preserved
- Familiar operational model

**Cons:**

- Still requires 24/7 process for infrequent jobs
- Distributed lock complexity remains
- Cost not proportional to actual work

### Option C: External Cron Service (GitHub Actions, CloudWatch)

**Description:**

Use external cron triggers that invoke API endpoints or queue jobs.

**Pros:**

- Free tier available (GitHub Actions)
- Platform-agnostic approach
- No Railway-specific configuration

**Cons:**

- Additional security surface (exposed endpoints)
- Cross-service coordination complexity
- Rate limiting and retry logic needed
- Monitoring fragmentation

## Implementation

### Removed Components

Files deleted from worker repository:

```
src/cmd/scheduler/main.go
src/internal/app/bootstrap/scheduler.go
src/internal/app/container_scheduler.go
src/internal/domain/analysis/autorefresh.go
src/internal/domain/analysis/decay.go
src/internal/handler/scheduler/autorefresh.go
src/internal/infra/scheduler/cron.go
src/internal/infra/scheduler/lock.go
src/internal/usecase/autorefresh/refresh.go
```

### New Binary Structure

```
cmd/
├── analyzer/           # Queue consumer (River)
├── spec-generator/     # Queue consumer (River)
├── retention-cleanup/  # Cron binary (Railway)
└── enqueue/            # Manual utility
```

### Infrastructure Configuration

```
infra/
├── analyzer/
│   ├── Dockerfile
│   └── railway.json
├── spec-generator/
│   ├── Dockerfile
│   └── railway.json
└── retention-cleanup/
    ├── Dockerfile
    └── railway.json
```

### Deployment Comparison

| Aspect            | Before (Scheduler)           | After (Railway Cron)           |
| ----------------- | ---------------------------- | ------------------------------ |
| Running cost      | 24/7 (even when idle)        | Per-execution only             |
| Distributed lock  | Required (PostgreSQL)        | Not needed (Railway manages)   |
| Scaling           | Fixed single instance        | Per-job independent            |
| Failure isolation | All jobs fail together       | Per-job isolation              |
| Configuration     | Environment variables + code | Railway IaC per service        |
| Monitoring        | Single service metrics       | Per-service metrics in Railway |

## Consequences

### Positive

**Cost Optimization:**

- Eliminates 24/7 running cost for infrequent cron jobs
- Pay only for actual execution time
- Per-job cost visibility enables optimization

**Operational Simplicity:**

- No distributed lock to manage or debug
- Each cron job is a simple run-to-completion binary
- Railway handles scheduling, retries, and single execution

**Deployment Independence:**

- Each cron job can be deployed independently
- Different schedules don't require code changes
- IaC-based configuration (infrastructure as code)

**Failure Isolation:**

- One failing cron job doesn't affect others
- Clear per-job logs and metrics
- Independent retry policies

### Negative

**Platform Dependency:**

- Tied to Railway's cron implementation
- Migration requires reconfiguring all cron jobs
- Railway-specific IaC format

**Cold Start Latency:**

- Each execution starts a new container
- Not suitable for sub-minute intervals
- Initial connection setup overhead

**Configuration Sprawl:**

- Multiple railway.json files to maintain
- Sync between Dockerfile and railway.json required
- More files in infra repository

### Superseded ADRs

| ADR                                                               | Status               | Notes                                                             |
| ----------------------------------------------------------------- | -------------------- | ----------------------------------------------------------------- |
| [Worker ADR-01](/en/adr/worker/01-scheduled-recollection.md)      | Superseded           | Auto-refresh scheduler removed                                    |
| [Worker ADR-05](/en/adr/worker/05-worker-scheduler-separation.md) | Partially Superseded | Scheduler no longer exists; binary separation pattern still valid |

### Related Updates Required

| Document                                                  | Update Needed                           |
| --------------------------------------------------------- | --------------------------------------- |
| [ADR-04](/en/adr/04-queue-based-async-processing.md)      | Add note about Railway Cron alternative |
| [ADR-12](/en/adr/12-worker-centric-analysis-lifecycle.md) | Remove Scheduler references             |

## References

- Commit `c163239`: Remove Scheduler service
- Commit `f3fae45`: Separate worker binaries
- Commit `6e03a7f`: Add retention-cleanup bootstrap
- [Railway Cron Documentation](https://docs.railway.app/reference/cron-jobs)
