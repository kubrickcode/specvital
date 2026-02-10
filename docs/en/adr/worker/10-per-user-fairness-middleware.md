---
title: Per-User Fairness Middleware
description: ADR for implementing per-user concurrent job limiting via River WorkerMiddleware
---

# ADR-10: Per-User Fairness Middleware

> 🇰🇷 [한국어 버전](/ko/adr/worker/10-per-user-fairness-middleware.md)

| Date       | Author       | Repos  |
| ---------- | ------------ | ------ |
| 2026-02-02 | @KubrickCode | worker |

## Context

### Problem Statement

Without per-user limits, a single user submitting many analysis requests monopolizes queue workers, causing unfair resource distribution:

| Scenario                    | Impact                                          |
| --------------------------- | ----------------------------------------------- |
| Free user submits 10 jobs   | All 5 workers occupied; Pro users wait in queue |
| Single user mass submission | Other users experience degraded service         |
| No tier differentiation     | Paid users receive no priority over free tier   |

### Requirements

| Requirement       | Description                                          |
| ----------------- | ---------------------------------------------------- |
| Per-User Limiting | Limit concurrent jobs per user, not globally         |
| Tier-Based Quotas | Higher tiers get more concurrent slots               |
| Non-Destructive   | Jobs delayed, not rejected; all work eventually runs |
| Low Overhead      | Minimal latency impact on job execution              |
| Graceful Handling | Prevent thundering herd when snoozed jobs wake       |

## Decision

**Implement per-user concurrent job limits via River WorkerMiddleware with tier-based quotas.**

### Tier Limits

| Tier       | Concurrent Jobs |
| ---------- | --------------- |
| Free       | 1               |
| Pro        | 3               |
| Pro Plus   | 3               |
| Enterprise | 5               |

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Fairness Middleware Flow                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Job Received                                                    │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────┐                                            │
│  │ UserJobExtractor │  Extract userID from job args (JSON)      │
│  └────────┬────────┘                                            │
│           │                                                      │
│           ▼                                                      │
│  ┌─────────────────┐                                            │
│  │  TierResolver   │  Query DB for user's subscription tier     │
│  │  (DB lookup)    │  Default to Free if not found              │
│  └────────┬────────┘                                            │
│           │                                                      │
│           ▼                                                      │
│  ┌─────────────────┐                                            │
│  │ PerUserLimiter  │  TryAcquire(userID, tier, jobID)          │
│  │  (in-memory)    │                                            │
│  └────────┬────────┘                                            │
│           │                                                      │
│     ┌─────┴─────┐                                               │
│     │           │                                                │
│  acquired    rejected                                            │
│     │           │                                                │
│     ▼           ▼                                                │
│  Execute    JobSnooze(30s + jitter)                             │
│  Worker     Return for retry                                     │
│     │                                                            │
│     ▼                                                            │
│  defer Release(userID, jobID)                                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Key Design Decisions

| Decision                    | Rationale                                                |
| --------------------------- | -------------------------------------------------------- |
| WorkerMiddleware over Hook  | Hooks cannot return JobSnooze; only middleware can       |
| In-memory limiter           | Per-instance state acceptable for single deployment      |
| Tier from DB (not job args) | Web layer doesn't include tier; maintains loose coupling |
| 30s + 10s jitter snooze     | Prevents thundering herd on wake                         |
| Idempotent acquire/release  | Same jobID cannot double-count slots                     |
| System job bypass           | Empty userID bypasses limits (scheduled jobs)            |

## Options Considered

### Option A: Per-User WorkerMiddleware Limiting (Selected)

Implement River WorkerMiddleware that tracks per-user concurrent job counts in memory. Before job execution, check tier-based limit. If exceeded, return JobSnooze.

**Pros:**

- Precise concurrent control per user
- Tier-aware fairness differentiation
- Jobs snoozed, not rejected
- Jitter prevents thundering herd

**Cons:**

- Per-instance state (not distributed)
- DB lookup for tier resolution

### Option B: Priority Field in Single Queue

Add priority field based on tier. Higher priority processed first.

**Rejected:** Priority affects ordering, not concurrency. Single user still monopolizes workers.

### Option C: Dedicated Worker Pools Per Tier

Separate queues with dedicated workers per tier.

**Rejected:** Resource inefficiency when pools unevenly loaded. Single user still monopolizes within-tier pool.

### Option D: River Hook Approach

Use HookWorkBegin to check limits.

**Rejected:** Hooks cannot return JobSnooze. Technical limitation, not preference.

## Consequences

**Positive:**

- Fair resource distribution across users
- Clear tier value proposition (1 vs 3 vs 5 slots)
- Non-destructive limiting (snooze, not reject)
- Consistent with existing semaphore pattern (ADR-06)

**Negative:**

- Per-instance state limits horizontal scaling
- Additional DB query per job for tier lookup
- 30s+ latency for over-limit users

**Operational:**

- Monitor snooze rate per tier
- Expose tier limits via environment variables
- Document distributed limiter as future work

## Configuration

```
FAIRNESS_ENABLED=true
FAIRNESS_FREE_LIMIT=1
FAIRNESS_PRO_LIMIT=3
FAIRNESS_ENTERPRISE_LIMIT=5
FAIRNESS_SNOOZE_DURATION=30s
FAIRNESS_SNOOZE_JITTER=10s
```

## References

- [ADR-21: Quota Reservation](/en/adr/21-quota-reservation.md) - Request-level quota protection
- [Worker ADR-06: Semaphore Clone Concurrency](/en/adr/worker/06-semaphore-clone-concurrency.md) - Similar in-memory limiter pattern
- [River WorkerMiddleware Documentation](https://riverqueue.com/docs/middleware)
- [GitHub Commits: 620849f, 527c1ae](https://github.com/specvital/worker)
