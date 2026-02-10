---
title: Phase 1 V1 Quality Enhancement Architecture
description: ADR for post-processing pipeline to improve classification quality in the V1 architecture
---

# ADR-09: Phase 1 V1 Quality Enhancement Architecture

> 🇰🇷 [한국어 버전](/ko/adr/worker/09-phase1-v1-quality-enhancement.md)

| Date       | Author       | Repos  |
| ---------- | ------------ | ------ |
| 2026-02-04 | @KubrickCode | worker |

## Context

### Problem Statement

The AI-based SpecView generation pipeline (ADR-14) uses Gemini for test classification, but raw LLM output exhibits quality issues that degrade user experience:

| Quality Issue              | Description                                                     | User Impact                    |
| -------------------------- | --------------------------------------------------------------- | ------------------------------ |
| Domain Naming Variance     | Same concept appears as "Auth", "Authentication", "AuthService" | Fragmented domain groupings    |
| Orphaned Tests             | Tests returned without domain assignment                        | Missing specifications         |
| Abbreviation Inconsistency | "db", "DB", "Database" treated as different domains             | Duplicate domains in output    |
| Structural Errors          | Occasional malformed JSON responses                             | Pipeline failures              |
| Quality Blind Spots        | No metrics on classification quality                            | Unable to measure improvements |

### Experimental Alternatives

Two alternative architectures were explored before returning to enhanced V1:

| Architecture        | Approach                                              | Outcome                                                           |
| ------------------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| V2 Two-Stage        | Separate domain extraction followed by classification | Abandoned - complexity overhead without proportional quality gain |
| V3 Sequential Batch | Process batches sequentially with anchor propagation  | Abandoned - latency increase outweighed consistency benefits      |

Both were disabled via feature flags (`SPECVIEW_PHASE1_V2=false`, `SPECVIEW_PHASE1_V3=false`) and code has been removed.

### Requirements

| Requirement         | Description                                      |
| ------------------- | ------------------------------------------------ |
| Backward Compatible | Must layer on existing ADR-14 pipeline           |
| Low Latency Impact  | Post-processing overhead < 100ms per chunk       |
| Observable          | Provide metrics for quality monitoring           |
| Deterministic       | Same input produces consistent domain assignment |

## Decision

**Enhance the V1 single-pass classification architecture with a post-processing pipeline that validates, normalizes, and recovers classification results.**

### Post-Processing Pipeline

```
LLM Classification Output
         │
         ▼
┌─────────────────────────────────────────────────────┐
│              Phase1PostProcessor                     │
├─────────────────────────────────────────────────────┤
│  1. JSON Validation        → Reject malformed       │
│  2. Domain Normalization   → Merge similar names    │
│  3. Abbreviation Expansion → auth → Authentication  │
│  4. Orphaned Detection     → Flag unclassified      │
│  5. Path-based Fallback    → Derive from file path  │
└─────────────────────────────────────────────────────┘
         │
         ▼
    Normalized Classification Result
         │
         ▼
┌─────────────────────────────────────────────────────┐
│           Quality Metrics Collector                  │
├─────────────────────────────────────────────────────┤
│  - Orphaned test count                              │
│  - Normalization frequency                          │
│  - Fallback usage rate                              │
│  - Domain distribution                              │
└─────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component                  | Responsibility                                       | Trigger                       |
| -------------------------- | ---------------------------------------------------- | ----------------------------- |
| Phase1PostProcessor        | Validate and normalize LLM classification output     | Every classification response |
| Domain Normalization       | Merge semantically equivalent domain names           | On domain name extraction     |
| Domain Abbreviation Expand | Expand common abbreviations to full names            | On domain name extraction     |
| Path-based Fallback        | Derive domain from test file path for orphaned tests | When test has no domain       |
| Orphaned Test Detection    | Identify and flag tests that failed classification   | After normalization           |
| Quality Metrics Collector  | Track classification quality metrics                 | After each batch              |

### Normalization Rules

**Domain Name Merging:**

| Variant                       | Normalized To  |
| ----------------------------- | -------------- |
| Auth, AuthService, AuthModule | Authentication |
| User, UserService, UserMgmt   | UserManagement |
| Pay, Payments, PaymentService | Payment        |
| DB, Database, DataAccess      | Database       |

**Abbreviation Expansion:**

| Abbreviation | Expanded       |
| ------------ | -------------- |
| auth         | Authentication |
| db           | Database       |
| ui           | UserInterface  |
| api          | API            |
| http         | HTTP           |

**Path-based Fallback Strategy:**

```
test/services/payment/checkout_test.go
         │
         ▼
Extract: "payment" from path
         │
         ▼
Expand:  "Payment" domain
```

## Options Considered

### Option A: V1 with Post-Processing Enhancement (Selected)

Add validation, normalization, and fallback layers after LLM response while maintaining the single-pass classification architecture from ADR-14.

| Aspect              | Assessment                                       |
| ------------------- | ------------------------------------------------ |
| Architecture Impact | Additive - no changes to core LLM pipeline       |
| Latency             | +10-50ms per chunk (post-processing)             |
| Quality Improvement | Moderate - addresses naming consistency          |
| Implementation      | 6 focused components with clear responsibilities |
| Rollback            | Feature-flaggable per component                  |

**Selection Rationale:**

- V2/V3 experiments demonstrated that additional LLM calls did not proportionally improve quality
- Post-processing addresses observed quality issues without API cost increase
- Maintains proven reliability characteristics of V1 pipeline
- Enables incremental improvement without architectural risk

### Option B: V2 Two-Stage Taxonomy Architecture (Abandoned)

Separate domain extraction into a dedicated LLM call before classification.

| Aspect       | Assessment                              |
| ------------ | --------------------------------------- |
| API Calls    | 2x (domain extraction + classification) |
| Complexity   | High - two-phase coordination           |
| Quality Gain | Minimal in testing                      |
| Cost Impact  | 2x Gemini API costs for Phase 1         |
| Status       | Abandoned - `SPECVIEW_PHASE1_V2=false`  |

**Rejection Rationale:**

- Empirical testing showed domain consistency did not improve proportionally
- Doubled API costs without user-visible quality improvement
- Added failure modes and debugging complexity

### Option C: V3 Sequential Batch Architecture (Abandoned)

Process test batches sequentially with explicit context carryover between batches.

| Aspect       | Assessment                                               |
| ------------ | -------------------------------------------------------- |
| Latency      | Increased - sequential processing eliminates parallelism |
| Complexity   | High - batch ordering and state management               |
| Quality Gain | Marginal consistency improvement                         |
| Status       | Abandoned - `SPECVIEW_PHASE1_V3=false`                   |

**Rejection Rationale:**

- Sequential processing significantly increased total processing time
- Anchor propagation in V1 already provides cross-chunk consistency
- Complexity not justified by quality improvement

## Consequences

### Positive

| Area                    | Benefit                                                               |
| ----------------------- | --------------------------------------------------------------------- |
| Domain Consistency      | Normalization eliminates duplicate domains from naming variance       |
| Test Coverage           | Path-based fallback recovers orphaned tests                           |
| Observability           | Metrics collector enables quality monitoring and regression detection |
| Latency                 | No additional API calls - post-processing is local computation        |
| Maintainability         | Each enhancement is isolated and independently testable               |
| Architecture Simplicity | Avoided V2/V3 complexity by improving existing pipeline               |

### Negative

| Area                   | Trade-off                                                        |
| ---------------------- | ---------------------------------------------------------------- |
| Quality Ceiling        | Single-pass LLM classification limits maximum achievable quality |
| Normalization Rules    | Manual curation required for domain synonym mappings             |
| Path Fallback Accuracy | File path may not accurately reflect business domain             |
| Metric Overhead        | Quality metrics add minor processing and storage overhead        |

### Technical Implications

| Aspect            | Implication                                                 |
| ----------------- | ----------------------------------------------------------- |
| Pipeline Position | PostProcessor runs synchronously after LLM response parsing |
| Error Handling    | Normalization failures log warning but do not fail pipeline |
| Metrics Storage   | Quality metrics logged via structured logging (no database) |
| Configuration     | Normalization rules configurable without code change        |
| Testing           | Golden snapshot tests for normalization behavior            |

## References

- [ADR-14: AI-Based Spec Document Generation Pipeline](/en/adr/14-ai-spec-generation-pipeline.md)
- [ADR-08: SpecView Worker Binary Separation](/en/adr/worker/08-specview-worker-separation.md)
- [Core ADR-16: Domain Hints Extraction System](/en/adr/core/16-domain-hints-extraction.md)
- [Related commits](https://github.com/specvital/worker/commits/main) - Phase1PostProcessor, MetricsCollector implementation
