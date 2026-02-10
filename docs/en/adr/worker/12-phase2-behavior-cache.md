---
title: Phase 2 Behavior Cache
description: ADR for caching Phase 2 AI-generated test behavior descriptions to reduce API costs
---

# ADR-12: Phase 2 Behavior Cache

> 🇰🇷 [한국어 버전](/ko/adr/worker/12-phase2-behavior-cache.md)

| Date       | Author       | Repos         |
| ---------- | ------------ | ------------- |
| 2026-01-24 | @KubrickCode | worker, infra |

## Context

### Phase 2 Cost Problem

The AI-Based Spec Document Generation Pipeline ([ADR-14](/en/adr/14-ai-spec-generation-pipeline.md)) uses a two-phase approach:

| Phase   | Model                 | Cost     | Purpose                         |
| ------- | --------------------- | -------- | ------------------------------- |
| Phase 1 | gemini-2.5-flash      | $0.30/1M | Test classification by domain   |
| Phase 2 | gemini-2.5-flash-lite | $0.10/1M | Test name → behavior conversion |

While Phase 1 results are cached at the document level via `content_hash`, Phase 2 involves per-test AI calls that are expensive when:

- Same test file is analyzed across multiple commits
- Similar test names appear across different repositories
- Re-analysis triggered by parser version updates or user requests

### Caching Opportunity

Test behavior descriptions have high cache reusability:

| Scenario                     | Example                                                         |
| ---------------------------- | --------------------------------------------------------------- |
| Same test, different commits | `test_user_login` in commit A and B produces identical behavior |
| Cross-repository similarity  | `testAuthentication` behavior is language/framework agnostic    |
| Re-analysis                  | Parser upgrade doesn't change behavior semantics                |

### Cache Key Design Challenge

**Problem**: What uniquely identifies a test behavior?

| Approach              | Pros            | Cons                                        |
| --------------------- | --------------- | ------------------------------------------- |
| Test name only        | Maximum reuse   | Ignores context (same name, different test) |
| Test name + file path | Context-aware   | Path changes invalidate cache               |
| Content hash of test  | Exact match     | No reuse across similar tests               |
| Semantic fingerprint  | Captures intent | Complex, requires additional AI call        |

## Decision

**Implement a PostgreSQL-backed behavior cache using composite key `(test_name_hash, language, model_id)` with TTL-based expiration.**

### Table Schema

```sql
CREATE TABLE behavior_caches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_name_hash TEXT NOT NULL,         -- SHA-256 of normalized test name
    language VARCHAR(10) NOT NULL,        -- en, ko, etc.
    model_id VARCHAR(100) NOT NULL,       -- gemini-2.5-flash-lite
    behavior_description TEXT NOT NULL,   -- Cached AI output
    confidence DECIMAL(3,2) NOT NULL,     -- 0.00-1.00
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,      -- TTL expiration
    hit_count INTEGER DEFAULT 0,          -- Usage tracking

    CONSTRAINT behavior_caches_unique
        UNIQUE (test_name_hash, language, model_id)
);

-- Index for lookup
CREATE INDEX idx_behavior_caches_lookup
    ON behavior_caches (test_name_hash, language, model_id)
    WHERE expires_at > NOW();

-- Index for cleanup
CREATE INDEX idx_behavior_caches_expiry
    ON behavior_caches (expires_at);
```

### Cache Key Strategy

```
test_name_hash = SHA256(normalize(test_name))

normalize(test_name):
  1. Lowercase conversion
  2. Remove common prefixes (test_, it_, describe_, should_)
  3. Strip special characters and numbers
  4. Collapse whitespace
```

**Examples**:

| Original Test Name                 | Normalized         | Hash (truncated) |
| ---------------------------------- | ------------------ | ---------------- |
| `test_user_can_login`              | `user can login`   | `a3f2...`        |
| `TestUserCanLogin`                 | `user can login`   | `a3f2...`        |
| `it('should allow user to login')` | `allow user login` | `b7c1...`        |
| `describe('User Login')`           | `user login`       | `c9e4...`        |

### Cache Lookup Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Phase 2 Processing                            │
├─────────────────────────────────────────────────────────────────┤
│  For each test in feature:                                       │
│                                                                  │
│  1. Compute test_name_hash                                       │
│                                                                  │
│  2. Cache lookup:                                                │
│     SELECT behavior_description, confidence                      │
│     FROM behavior_caches                                         │
│     WHERE test_name_hash = ? AND language = ? AND model_id = ?   │
│       AND expires_at > NOW()                                     │
│                                                                  │
│  3. If HIT:                                                      │
│     - Increment hit_count                                        │
│     - Return cached behavior (skip AI call)                      │
│     - No quota consumption                                       │
│                                                                  │
│  4. If MISS:                                                     │
│     - Call Gemini API                                            │
│     - Store result with TTL                                      │
│     - Consume quota                                              │
└─────────────────────────────────────────────────────────────────┘
```

### TTL Configuration

| Tier       | TTL      | Rationale                   |
| ---------- | -------- | --------------------------- |
| Free       | 7 days   | Limited storage, high churn |
| Pro        | 30 days  | Standard business retention |
| Pro Plus   | 90 days  | Extended caching value      |
| Enterprise | 180 days | Maximum cost optimization   |

### Integration with Quota System

Per [ADR-13](/en/adr/13-billing-quota-architecture.md), cache hits do not consume quota:

```go
type Phase2Result struct {
    Behavior   string
    Confidence float64
    FromCache  bool  // If true, quota not consumed
}

// In usage tracking
if !result.FromCache {
    quotaService.RecordUsage(ctx, userID, QuotaTypeSpecView)
}
```

## Options Considered

### Option A: PostgreSQL Table (Selected)

Database-backed cache with TTL and hit tracking.

**Pros**:

- Integrated with existing PostgreSQL infrastructure
- Queryable for analytics (hit rates, popular tests)
- Automatic cleanup via scheduled job
- Transactional consistency with spec document writes

**Cons**:

- Database load for high-volume lookups
- Storage costs for large caches

### Option B: Redis Cache

In-memory cache with automatic expiration.

**Pros**:

- Sub-millisecond lookups
- Native TTL support
- Reduced database load

**Cons**:

- Additional infrastructure (not currently in stack)
- Cache loss on Redis restart
- Memory cost scaling with cache size

### Option C: Document-Level Caching Only

Rely on existing `spec_documents.content_hash` caching.

**Pros**:

- No new infrastructure
- Already implemented

**Cons**:

- No reuse for similar tests across repositories
- Full Phase 2 re-run on any test change
- Misses cross-document optimization opportunity

### Option D: No Additional Caching

Accept Phase 2 costs as operational expense.

**Pros**:

- Simplest implementation
- No cache invalidation complexity

**Cons**:

- Higher API costs at scale
- Slower response times for repeated tests
- Poor cost efficiency for high-volume users

## Consequences

### Positive

| Area           | Benefit                                          |
| -------------- | ------------------------------------------------ |
| Cost Reduction | 40-60% Phase 2 API cost savings (estimated)      |
| Response Time  | Cache hits avoid 1-2s AI latency per test        |
| Quota Fairness | Cache hits don't consume user quota              |
| Analytics      | Hit rate metrics enable cache tuning             |
| Cross-Repo     | Similar tests across repos share cached behavior |

### Negative

| Area                 | Trade-off                                        |
| -------------------- | ------------------------------------------------ |
| Storage              | Cache table grows with unique test diversity     |
| Stale Data Risk      | Cached behavior may not reflect model updates    |
| Normalization Errors | Aggressive normalization may cause false matches |
| Cleanup Overhead     | Scheduled job required for TTL enforcement       |

### Technical Notes

- **Cache warming**: Not implemented; cache builds organically
- **Invalidation**: Model version change invalidates via `model_id` in key
- **Conflict resolution**: First writer wins; concurrent writes are rare
- **Monitoring**: Log cache hit/miss ratio per analysis

## Configuration

```
BEHAVIOR_CACHE_ENABLED=true
BEHAVIOR_CACHE_DEFAULT_TTL=30d
BEHAVIOR_CACHE_CLEANUP_SCHEDULE=0 4 * * *  # 4 AM UTC daily
BEHAVIOR_CACHE_CLEANUP_BATCH_SIZE=5000
```

## References

- [ADR-14: AI-Based Spec Document Generation Pipeline](/en/adr/14-ai-spec-generation-pipeline.md) - Parent architecture
- [ADR-13: Billing and Quota Architecture](/en/adr/13-billing-quota-architecture.md) - Quota integration
- [ADR-18: GitHub API Cache Tables](/en/adr/18-github-api-cache-tables.md) - Similar caching pattern
- Commit: `8917156` (behavior_caches table) - 2026-01-24
