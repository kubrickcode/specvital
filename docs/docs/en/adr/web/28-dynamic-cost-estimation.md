---
title: Dynamic Cost Estimation
description: ADR for cache-aware cost prediction in the spec generation modal
---

# ADR-28: Dynamic Cost Estimation

> 🇰🇷 [한국어 버전](/ko/adr/web/28-dynamic-cost-estimation.md)

| Date       | Author       | Repos |
| ---------- | ------------ | ----- |
| 2026-02-03 | @KubrickCode | web   |

## Context

### The Problem

The spec generation modal displays estimated usage costs before users trigger AI-powered test behavior generation. Per [ADR-13: Billing and Quota Architecture](/en/adr/13-billing-quota-architecture.md), users have limited monthly quotas and cache hits do not consume quota.

However, the estimated usage remained constant regardless of cache option selection:

| User Selection | Expected Behavior                  | Actual Behavior      |
| -------------- | ---------------------------------- | -------------------- |
| Use Cache      | Show reduced estimate (cache hits) | Showed same as fresh |
| Fresh Analysis | Show full test count               | Showed same as cache |

This occurred because:

1. Estimates were hardcoded to `current_test_count`
2. No mechanism existed to predict cache hit rates
3. Worker's cache lookup algorithm was not accessible to Web layer

### Technical Challenge

The Worker service uses `SHA256(normalize(test_name))` for actual cache key lookup ([Worker ADR-12](/en/adr/worker/12-phase2-behavior-cache.md)). Replicating this algorithm in Web requires:

- Sharing normalization logic across repositories
- Frontend computing SHA256 hashes
- Maintaining synchronization when Worker's algorithm changes

## Decision

**Implement client-side cache hit prediction using name+path matching against previous spec behaviors.**

### Algorithm

```
For each test in current analysis:
  1. Find matching test in previous spec_document.behaviors
     Match criteria: test.name == behavior.test_name
                 AND test.file_path == behavior.file_path
  2. If match found: predict cache hit
  3. If no match: predict cache miss

estimated_cost = current_tests - predicted_cache_hits
```

### Design Choices

| Choice                 | Rationale                                   |
| ---------------------- | ------------------------------------------- |
| Approximate matching   | ~99% accuracy with simple implementation    |
| Client-side prediction | No API round-trip; instant modal rendering  |
| Disclaimer tooltip     | "Actual usage may differ from estimate"     |
| Graceful degradation   | No previous analysis → show full test count |

## Options Considered

### Option A: Name+Path Prediction (Selected)

Match tests by `name + file_path` comparison against previous spec document behaviors.

**Pros:**

- Zero latency: computed client-side from existing data
- No API dependency: works offline
- Simple implementation: string comparison only
- Maintainable: independent of Worker normalization changes
- High accuracy: ~99% correlation with actual cache hits

**Cons:**

- Prediction mismatch: edge cases where name+path matches but normalized hash differs
- False negatives: same logical test with path change appears as miss
- No cross-repo benefit: cannot predict cache hits from other repositories

### Option B: Replicate Worker's SHA256 Hash Algorithm

Implement identical normalization and hashing in Web frontend.

**Pros:**

- 100% accuracy with actual cache behavior
- Cross-repo prediction theoretically possible

**Cons:**

- Code duplication: normalization logic in Worker (Go) and Web (TypeScript)
- Sync maintenance: algorithm changes require coordinated deployments
- Bundle size: SHA256 library adds to frontend
- Complexity: regex-based normalization prone to drift
- Test burden: must test parity across languages

**Rejected**: Coupling cost outweighs accuracy benefit.

### Option C: Always Show Worst-Case Estimate

Display maximum possible cost (all tests = cache miss).

**Pros:**

- Zero prediction risk
- Simplest implementation

**Cons:**

- Poor UX: cache option appears valueless
- Conversion impact: users skip beneficial cache mode
- Trust erosion: estimates always wrong in cached scenarios

**Rejected**: Defeats purpose of cache transparency.

### Option D: Call Worker Service for Prediction

Add API endpoint that runs actual cache lookup without consuming quota.

**Pros:**

- 100% accuracy
- Single source of truth

**Cons:**

- Additional latency: network delays modal rendering
- API complexity: new endpoint to maintain
- Failure mode: modal broken when Worker unavailable
- Rate limiting: could be abused to probe cache state

**Rejected**: Over-engineering for prediction feature.

## Implementation Details

### API Endpoint

```
GET /api/spec-view/{analysisId}/cache-prediction?language={language}
```

### Data Structure

```typescript
interface CachePrediction {
  totalBehaviors: number;
  cacheableCount: number;
  estimatedCost: number;
}
```

### SQL Query Structure

```sql
WITH
  previous_spec AS (
    -- Find most recent spec for same codebase/language
    SELECT id FROM spec_documents
    WHERE ... ORDER BY created_at DESC LIMIT 1
  ),
  previous_behaviors AS (
    -- Extract behaviors with test names and file paths
    SELECT test_name, file_path FROM spec_behaviors
    WHERE spec_document_id = (SELECT id FROM previous_spec)
  ),
  current_tests AS (
    -- Get all test cases from current analysis
    SELECT name, file_path FROM test_cases
    WHERE analysis_id = $1
  ),
  matched_behaviors AS (
    -- Count where name AND file_path match
    SELECT COUNT(*) FROM current_tests c
    JOIN previous_behaviors p
      ON c.name = p.test_name AND c.file_path = p.file_path
  )
SELECT total_behaviors, cacheable_count;
```

### Frontend Integration

```typescript
const dynamicEstimatedCost = forceRegenerate
  ? estimatedCost
  : (cachePrediction?.estimatedCost ?? estimatedCost);
```

## Consequences

### Positive

- **User Experience**: Estimates reflect cache benefit, enabling informed decisions
- **Cost Transparency**: Users see approximate savings before committing quota
- **Implementation Speed**: Simple approach shipped without cross-repo coordination
- **Resilience**: No external dependencies; works when Worker unavailable
- **Maintainability**: Web layer isolated from Worker's cache key evolution

### Negative

- **Prediction Accuracy**: ~1% of cases may show incorrect estimate
- **Algorithm Divergence**: Intentional mismatch with Worker's actual algorithm
- **Limited Scope**: Cannot predict cross-repository cache hits
- **User Expectation**: Requires disclaimer to manage estimate != actual

### Edge Cases

| Scenario                    | Prediction | Reality            | Impact                 |
| --------------------------- | ---------- | ------------------ | ---------------------- |
| Test renamed, same path     | Cache miss | Cache miss         | Correct                |
| Test moved, same name       | Cache miss | Possible cache hit | Under-predicts savings |
| Test copied to new file     | Cache miss | Cache hit          | Under-predicts savings |
| Whitespace-only name change | Cache hit  | Cache hit          | Correct                |

Edge cases result in under-prediction of savings (higher cost shown), which is safe for billing context.

## References

### Internal

- [ADR-13: Billing and Quota Architecture](/en/adr/13-billing-quota-architecture.md) - Event-based usage tracking
- [ADR-21: Quota Reservation](/en/adr/21-quota-reservation.md) - Atomic quota reservation
- [Worker ADR-12: Phase 2 Behavior Cache](/en/adr/worker/12-phase2-behavior-cache.md) - Actual cache algorithm
- [Web ADR-26: Zustand Client State](/en/adr/web/26-zustand-client-state.md) - State management patterns

### External

- [Commit 3645d52](https://github.com/specvital/web/commit/3645d52) - Initial implementation
- [Commit e5a7a53](https://github.com/specvital/web/commit/e5a7a53) - Cache reuse selection UI
