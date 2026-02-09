---
title: Commit-Based Analysis Browsing
description: ADR for URL-based commit selection to enable historical analysis browsing
---

# ADR-27: Commit-Based Analysis Browsing

> 🇰🇷 [한국어 버전](/ko/adr/web/27-commit-based-analysis-browsing.md)

| Date       | Author       | Repos |
| ---------- | ------------ | ----- |
| 2026-02-03 | @KubrickCode | web   |

## Context

### The Problem

Users could only view the latest analysis results for a repository. However, the database stores multiple analyses tied to specific commit SHAs. After force-push operations (`git reset --hard && git push -f`), users saw "Up to date" without context, and previously analyzed commits became inaccessible.

### User Needs

| Need                     | Description                                                |
| ------------------------ | ---------------------------------------------------------- |
| Historical browsing      | View analysis results from any previously analyzed commit  |
| Comparison over time     | Understand how test specifications evolved                 |
| Debug regressions        | Identify when specific test specification changes occurred |
| Shareable analysis state | Send colleagues a link to a specific commit's analysis     |

### Existing Constraints

| Constraint         | Source     | Implication                                       |
| ------------------ | ---------- | ------------------------------------------------- |
| nuqs for URL state | Web ADR-16 | Commit selection should use URL parameters        |
| Zustand scope      | Web ADR-26 | Ephemeral state only; not suitable for navigation |
| TanStack Query     | Web ADR-04 | Analysis data fetching pattern established        |
| Next.js App Router | Web ADR-07 | SSR considerations for URL state                  |

## Decision

**Adopt URL-based commit selection using nuqs with a CommitSelector UI component.**

Core principles:

1. **URL as Source of Truth**: Commit selection persists in URL (`?commit=abc1234`)
2. **Progressive Enhancement**: Default to latest commit when no parameter specified
3. **nuqs Integration**: Use `parseAsString` with SHA format validation
4. **Compound State**: Combine commit selection with existing filters

### URL Structure

```
/repos/{owner}/{repo}/analysis?commit=abc1234
/repos/{owner}/{repo}/analysis?commit=abc1234&view=starred&q=auth
```

## Options Considered

### Option A: URL-Based Commit Selection via nuqs (Selected)

Store selected commit SHA in URL query parameter using the established nuqs pattern.

**Pros:**

- Shareable URLs with exact analysis state
- Browser history integration (back/forward navigation)
- SSR compatible for pre-fetching
- Consistent with existing filter patterns (ADR-16)
- Composable with other URL state

**Cons:**

- URL length increase (mitigated by 7-char SHA abbreviation)
- Invalid commit handling required
- Stale bookmarks when commits deleted

### Option B: Client-Side State via Zustand

Store selected commit in Zustand store as ephemeral client state.

**Pros:**

- Clean URLs
- Fast state changes
- Cross-component access

**Cons:**

- Not shareable
- Lost on page refresh
- No browser history support
- Violates ADR-26 scope (Zustand for ephemeral state only)

**Rejected**: Contradicts established patterns; poor UX for shareable links.

### Option C: Path Segment-Based Selection

Encode commit SHA in URL path (`/analysis/{commit}`).

**Pros:**

- RESTful semantics
- Cleaner URL hierarchy

**Cons:**

- Routing complexity
- Inconsistent with query-based filters
- Page transitions vs. component updates

**Rejected**: Higher implementation effort; breaks filter composition.

### Option D: Dropdown Only (No State Persistence)

UI component with local state only.

**Rejected**: Unacceptable UX; selection lost on any navigation.

## Implementation Details

### API Endpoints

| Endpoint                                       | Purpose                          |
| ---------------------------------------------- | -------------------------------- |
| `GET /api/analyze/{owner}/{repo}/history`      | List analyzed commits (max 50)   |
| `GET /api/analyze/{owner}/{repo}?commit={sha}` | Get analysis for specific commit |

### Data Structure

```typescript
interface AnalysisHistoryItem {
  id: string;
  commitSha: string;
  commitDate: string | null;
  completedAt: string;
  branchName: string | null;
  testCount: number;
  isHead: boolean;
}
```

### URL State Hook

```typescript
import { parseAsString, useQueryState } from "nuqs";

const commitParser = parseAsString.withDefault("");

export const useCommitSelect = () => {
  const [commit, setCommit] = useQueryState("commit", commitParser);
  const validated = /^[a-f0-9]{7,40}$/.test(commit ?? "") ? commit : null;
  return { commit: validated, setCommit } as const;
};
```

### CommitSelector Component

- Dropdown displaying abbreviated SHA (7 chars)
- `[HEAD]` badge for current HEAD commit
- Shows test count, branch name, date
- Falls back to static text when < 2 history entries

### TanStack Query Integration

```typescript
export const useAnalysisData = (owner: string, repo: string) => {
  const { commit } = useCommitSelect();

  return useQuery({
    queryKey: ["analysis", owner, repo, commit],
    queryFn: () => fetchAnalysis(owner, repo, commit),
  });
};
```

## Consequences

### Positive

- **Shareable Analysis State**: URLs like `/analysis?commit=abc1234&view=starred` enable collaboration
- **Browser Navigation**: Back/forward buttons navigate commit history naturally
- **SSR Optimization**: Server can pre-fetch correct analysis from URL parameter
- **Pattern Consistency**: Extends nuqs patterns from ADR-16
- **Debuggability**: Support can request "send me your URL" for exact state

### Negative

- **URL Complexity**: Longer URLs with commit parameter (mitigated by SHA abbreviation)
- **Invalid State Handling**: Must gracefully handle deleted/unavailable commits
- **Query Invalidation**: Changing commit invalidates TanStack Query cache (expected behavior)

## References

### Internal

- [ADR-04: TanStack Query Selection](/en/adr/web/04-tanstack-query-selection.md) - Server state management
- [ADR-16: nuqs URL State Management](/en/adr/web/16-nuqs-url-state-management.md) - URL state patterns
- [ADR-26: Zustand Client State Management](/en/adr/web/26-zustand-client-state.md) - State management trifecta
- [ADR-12: Worker-Centric Analysis Lifecycle](/en/adr/12-worker-centric-analysis-lifecycle.md) - Analysis data model

### External

- [GitHub Issue #300](https://github.com/specvital/web/issues/300) - Feature request and design
- [Release v1.6.0](https://github.com/specvital/web/releases/tag/v1.6.0) - Initial implementation
