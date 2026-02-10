---
title: Zustand Client State Management
description: ADR for client-side state management using Zustand in the web platform
---

# ADR-26: Zustand Client State Management

> 🇰🇷 [한국어 버전](/ko/adr/web/26-zustand-client-state)

| Date       | Author       | Repos |
| ---------- | ------------ | ----- |
| 2026-01-23 | @KubrickCode | web   |

## Context

The Specvital web platform requires **ephemeral client-side state** that differs from existing state management solutions:

| State Category   | Managed By     | Examples                                         |
| ---------------- | -------------- | ------------------------------------------------ |
| **Server State** | TanStack Query | Repository data, analysis results, user profile  |
| **URL State**    | nuqs           | Dashboard filters, search queries, view mode     |
| **Client State** | _Needed_       | Background task tracking, session ephemeral data |

### Problem

Background operations (analysis, spec generation) had critical limitations with modal-based UI:

1. State lost when modals closed
2. Polling stopped on page navigation
3. No cross-page visibility of running tasks

User quote from Issue #240: _"Users cannot track spec generation progress after closing modals or navigating away."_

### Constraints

- **App Router RSC Boundary**: State management must respect client/server component boundary
- **Hydration Mismatch Risk**: SSR + client state requires careful handling
- **Bundle Size Sensitivity**: Dashboard performance is critical
- **Existing Stack Integration**: Must coexist with TanStack Query and nuqs

## Decision

Adopt **Zustand** for global client-side state management.

### Rationale

1. Completes state management trifecta: TanStack Query (server), nuqs (URL), Zustand (client)
2. Minimal bundle impact (~3 KB gzipped)
3. No Provider wrapper required
4. Built-in middleware: `persist`, `devtools`
5. Native `useSyncExternalStore` for React 18+ optimization
6. TypeScript-first with full type inference

## Options Considered

### Option A: Zustand (Selected)

**Characteristics**:

- ~3 KB gzipped bundle size
- No Provider wrapper required
- Hooks-based API with selective subscriptions

**Strengths**:

- Zero boilerplate, immediate productivity
- Works anywhere in component tree without setup
- Redux DevTools compatible
- Built-in persist middleware for storage sync

**Weaknesses**:

- Less structured than Redux for very large codebases
- Team must establish conventions (no enforced patterns)

### Option B: Redux Toolkit

**Characteristics**:

- ~10-15 KB gzipped
- Provider wrapper required
- RTK Query for server state

**Evaluation**: Rejected. RTK Query overlaps with existing TanStack Query. Additional boilerplate and provider requirements not justified for remaining client state needs.

### Option C: React Context + useReducer

**Characteristics**:

- 0 KB (built-in)
- Provider wrapper required
- Manual optimization needed

**Evaluation**: Rejected. Re-render issues with frequently-changing state. Context is designed for low-frequency updates, not dynamic state like task tracking.

### Option D: Jotai (Atomic State)

**Characteristics**:

- ~3.5 KB gzipped
- Atom-based composition
- Provider optional but recommended

**Evaluation**: Viable but not selected. Atom-based model requires different mental model. Zustand's store-based approach is more intuitive for simple task tracking needs.

## Implementation

### Store Structure

```typescript
type TaskStoreState = {
  tasks: Map<string, BackgroundTask>;
};

type TaskStoreActions = {
  addTask: (task: Omit<BackgroundTask, "createdAt">) => void;
  updateTask: (id: string, updates: Partial<BackgroundTask>) => void;
  removeTask: (id: string) => void;
  clearCompletedTasks: () => void;
};
```

### Key Files

| File                                             | Purpose                      |
| ------------------------------------------------ | ---------------------------- |
| `lib/background-tasks/task-store.ts`             | Main Zustand store           |
| `lib/background-tasks/hooks.ts`                  | Custom hooks for consumers   |
| `lib/background-tasks/task-store.spec.ts`        | Unit tests                   |
| `lib/background-tasks/components/task-badge.tsx` | Badge with active task count |

### Patterns Used

| Pattern              | Implementation                                      |
| -------------------- | --------------------------------------------------- |
| Singleton Store      | Single `useTaskStore` with `create()`               |
| Derived Selectors    | `useActiveTasks` computes filtered results          |
| Shallow Comparison   | `useShallow` prevents unnecessary re-renders        |
| Outside React Access | `useTaskStore.getState()` for non-component usage   |
| Manual Persistence   | Custom sessionStorage sync (not persist middleware) |

### Persistence

Custom sessionStorage implementation instead of Zustand's `persist` middleware:

- Storage key: `specvital:background-tasks`
- SSR-safe with `typeof window === "undefined"` check
- Silent error handling for corrupted storage or quota exceeded

## Consequences

### Positive

- **Minimal Integration Friction**: No Provider needed, works immediately
- **Developer Experience**: Single file per store with collocated state + actions
- **Performance**: Selective subscriptions prevent unnecessary re-renders
- **Debugging**: Redux DevTools compatibility

### Negative

- **Team Consistency Risk**: No enforced patterns
  - **Mitigation**: Document store organization conventions, code review enforcement
- **Limited Middleware Ecosystem**: Fewer community middlewares than Redux
  - **Mitigation**: Core needs (persist, devtools) are covered
- **Next.js SSR Considerations**: Store is module-level
  - **Mitigation**: Follow official Zustand Next.js guide

## References

### Internal

- [ADR-04: TanStack Query Selection](/en/adr/web/04-tanstack-query-selection) - Explicitly stated "No Global State Library"
- [ADR-16: nuqs URL State Management](/en/adr/web/16-nuqs-url-state-management) - URL/shareable state

### External

- [Zustand Official Documentation](https://zustand.docs.pmnd.rs/)
- [Zustand GitHub Repository](https://github.com/pmndrs/zustand)
- [Setup with Next.js - Zustand Guide](https://zustand.docs.pmnd.rs/guides/nextjs)

### Related Commits

- `8664cbc`: feat(background-tasks): add global task store with persistence
- `ecb4434`: feat(background-tasks): integrate Account Badge and Tasks Dropdown
- `fc99ce5`: feat(background-tasks): add Dashboard active tasks section
- `2e1e6df`: refactor(dashboard): migrate use-reanalyze hook to TanStack Query polling
