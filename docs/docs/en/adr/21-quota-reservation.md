---
title: Quota Reservation for Concurrent Request Handling
description: ADR for preventing race conditions in quota checks during concurrent requests
---

# ADR-21: Quota Reservation for Concurrent Request Handling

> 🇰🇷 [한국어 버전](/ko/adr/21-quota-reservation)

| Date       | Author     | Repos              |
| ---------- | ---------- | ------------------ |
| 2026-02-01 | @specvital | web, worker, infra |

## Context

### The Concurrency Gap in Event-Based Billing

[ADR-13 (Billing and Quota Architecture)](/en/adr/13-billing-quota-architecture) established event-based usage tracking where usage events are recorded at job completion. This model ensures billing accuracy (only successful operations are charged) and aligns with the cache-first architecture (cache hits generate no events).

However, a temporal gap exists between quota check (request submission) and quota consumption (job completion). During high concurrency, this gap creates a race condition:

**Race Condition Scenario:**

| Time | Action                  | User State                |
| ---- | ----------------------- | ------------------------- |
| T0   | User has 4998/5000 used | -                         |
| T1   | Request A checks quota  | 4998 < 5000 → Pass        |
| T2   | Request B checks quota  | 4998 < 5000 → Pass        |
| T3   | Request A completes     | 5008 used (exceeds limit) |
| T4   | Request B completes     | 5018 used                 |

**Consequences of Unaddressed Race Condition:**

- Server resources wasted processing over-quota jobs
- Billing discrepancies where actual usage exceeds limits
- Unfair system behavior depending on request timing

### Constraints

| Constraint                        | Rationale                                          |
| --------------------------------- | -------------------------------------------------- |
| Must block at Web level           | Worker re-check wastes queue and compute resources |
| PostgreSQL-only solution          | Align with ADR-04 (River uses PostgreSQL)          |
| No user-visible billing anomalies | Preserve user trust in billing accuracy            |
| Worker owns lifecycle             | ADR-12 establishes Worker as record creator        |

## Decision

**Adopt quota reservation pattern with atomic transactions for concurrent request handling.**

### Reservation Mechanism

Introduce `quota_reservations` table to track in-flight quota commitments:

| Column          | Type             | Purpose                        |
| --------------- | ---------------- | ------------------------------ |
| id              | uuid (PK)        | Primary key                    |
| user_id         | uuid (FK)        | Links to user account          |
| event_type      | usage_event_type | specview or analysis           |
| reserved_amount | int              | Anticipated quota consumption  |
| job_id          | bigint (UNIQUE)  | Links to River job for cleanup |
| expires_at      | timestamptz      | 1-hour TTL for orphan cleanup  |

### Quota Check Formula

```sql
used + reserved + requested_amount <= limit
```

Where:

- `used`: Sum of usage_events in current billing period
- `reserved`: Sum of active reservations (expires_at > NOW())
- `requested_amount`: Units required for current request

### Transaction Atomicity

Web layer executes in single PostgreSQL transaction:

1. Check quota (including active reservations)
2. Create reservation record
3. Insert job into River queue (InsertTx)

Atomicity ensures: if job insertion fails, reservation is not created (rollback).

### Reservation Lifecycle

```
Web: Check quota → Create reservation → InsertTx
                          ↓
Worker: Process job → Delete reservation (success or failure)
                          ↓
Cleanup: Expire orphaned reservations after 1 hour
```

## Options Considered

### Option A: Reservation Pattern (Selected)

**Mechanism**: Create reservation atomically with job insertion; Worker deletes on completion.

| Aspect              | Assessment                                    |
| ------------------- | --------------------------------------------- |
| Quota accuracy      | Guaranteed - no concurrent over-commitment    |
| User experience     | Transparent - reservations are internal state |
| Resource efficiency | High - over-quota blocked at Web              |
| Complexity          | Moderate - additional table and cleanup job   |

### Option B: Pre-Deduction with Refund

**Mechanism**: Deduct quota at submission; refund on job failure.

| Aspect              | Assessment                                              |
| ------------------- | ------------------------------------------------------- |
| Quota accuracy      | Guaranteed                                              |
| User experience     | **Poor** - users see inflated usage during processing   |
| Resource efficiency | High                                                    |
| Complexity          | **High** - multiple failure states require refund logic |

**Rejected**: Creates billing anxiety and support burden. Users observing dashboards during job processing would see usage that may be refunded, undermining trust.

### Option C: Worker Re-Check

**Mechanism**: Optimistic check at Web; authoritative check at Worker.

| Aspect              | Assessment                                           |
| ------------------- | ---------------------------------------------------- |
| Quota accuracy      | Guaranteed (at Worker level)                         |
| User experience     | **Poor** - job accepted then rejected asynchronously |
| Resource efficiency | **Poor** - over-quota jobs consume queue capacity    |
| Complexity          | Low                                                  |

**Rejected**: Violates requirement to block at Web level. Wastes queue and Worker resources on jobs that will be rejected.

## Consequences

### Positive

| Benefit                    | Impact                                        |
| -------------------------- | --------------------------------------------- |
| Accurate quota enforcement | Concurrent requests cannot exceed limits      |
| Resource efficiency        | No wasted Worker cycles on over-quota jobs    |
| User fairness              | Only completed work appears in billing        |
| Atomic consistency         | Job + reservation in single transaction       |
| Clean failure handling     | Reservation deleted regardless of job outcome |

### Negative

| Trade-off                                  | Mitigation                                        |
| ------------------------------------------ | ------------------------------------------------- |
| Additional schema complexity               | Standard pattern; single table with clear purpose |
| Query overhead for reservation aggregation | Index on (user_id, event_type, expires_at)        |
| Orphan cleanup requirement                 | Scheduled job with configurable 1-hour TTL        |
| Debugging complexity                       | Structured logging at reserve/release points      |

### Technical Implications

| Aspect              | Requirement                                                    |
| ------------------- | -------------------------------------------------------------- |
| Transaction scope   | InsertTx and CreateReservation share PostgreSQL transaction    |
| Index design        | Compound index for efficient reservation lookup and expiration |
| Cleanup schedule    | Cron job to delete reservations where expires_at < NOW()       |
| Monitoring          | Alert on orphan count (indicates Worker health issues)         |
| Worker modification | Delete reservation by job_id on job completion                 |

## References

- [ADR-13: Billing and Quota Architecture](/en/adr/13-billing-quota-architecture)
- [ADR-04: Queue-Based Async Processing](/en/adr/04-queue-based-async-processing)
- [ADR-12: Worker-Centric Analysis Lifecycle](/en/adr/12-worker-centric-analysis-lifecycle)
- [GitHub Issue #291: Quota Reservation System](https://github.com/specvital/web/issues/291)
- Commit `d3f15bf`: feat(db): add quota_reservations table for concurrent request handling
