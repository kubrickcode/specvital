# Frontend CLAUDE.md

Next.js 16 App Router with React 19, acting as Backend-for-Frontend.

## Framework Guidelines

@nextjs.md

## Architecture: Feature-Based BFF

```
app/[locale]/ → features/[name]/ → lib/api/
```

| Layer           | Responsibility                          |
| --------------- | --------------------------------------- |
| **app/**        | Routes, layouts, Server Components      |
| **features/**   | Domain modules (components, hooks, api) |
| **components/** | Shared UI (shadcn/ui, layout, theme)    |
| **lib/**        | Utilities, API client, styles           |

### Component Strategy

| Type                 | Use When                            | Example                        |
| -------------------- | ----------------------------------- | ------------------------------ |
| **Server Component** | Data fetching, no interactivity     | Page components, data displays |
| **Client Component** | useState, useEffect, event handlers | Forms, modals, toggles         |

```tsx
// Server Component (default)
export default async function Page() {
  const data = await fetchData();
  return <Display data={data} />;
}

// Client Component (explicit)
("use client");
export function InteractiveForm() {
  const [state, setState] = useState();
}
```

## Auto-Generated Files (NEVER modify)

| File                         | Source                        | Regenerate     |
| ---------------------------- | ----------------------------- | -------------- |
| `lib/api/generated-types.ts` | `../backend/api/openapi.yaml` | `just gen-api` |

## Common Workflows

### New Feature

Create `features/[name]/`:

```
features/[name]/
├── components/
│   ├── feature-component.tsx
│   └── index.ts
├── hooks/           # (optional)
├── api/             # (optional)
└── index.ts         # Barrel export
```

### New Page

1. Create `app/[locale]/[route]/page.tsx`
2. Add Server Component for data fetching
3. Use feature components for UI
4. Add i18n messages if needed

### i18n

Add messages to `messages/{en,ko}.json`, use `useTranslations`:

```tsx
const t = useTranslations("namespace");
return <p>{t("key")}</p>;
```

## API Client

```tsx
import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { AnalysisResponse } from "@/lib/api/types";

const response = await apiFetch(`/api/analyze/${owner}/${repo}`);
const data = await parseJsonResponse<AnalysisResponse>(response);
```

### Token Auto-Refresh

Client-side API calls auto-refresh on 401:

1. Access token expires (15min)
2. Interceptor calls `/api/auth/refresh`
3. New token pair set via cookies
4. Original request retried

Key file: `lib/api/token-refresh.ts`

## Environment Variables

| Variable              | Purpose     | Default               |
| --------------------- | ----------- | --------------------- |
| `NEXT_PUBLIC_API_URL` | Backend URL | http://localhost:8000 |

## React Compiler

React Compiler is enabled (`next.config.ts`: `reactCompiler: true`)

### Prohibited

- **NEVER** use `useMemo`, `useCallback`, `React.memo`
- React Compiler handles memoization automatically at build time

### Escape Hatch

Use `"use no memo"` directive only when compiler causes issues:

```tsx
"use no memo";

export function SpecialComponent() {
  // Excluded from React Compiler optimization
}
```

## React Query Caching Policy

**Default: No caching (`staleTime: 0`)**

- Always provide fresh data on mount/refetch
- **NEVER** set `staleTime` unless data is truly static (e.g., pricing info)
- If `staleTime` must be set, `gcTime` MUST equal or exceed `staleTime`

**Rationale**: Arbitrary caching causes hard-to-debug issues. The cost of an extra API call is far less than debugging cache-related bugs.

**Prohibited Patterns**:

- `isFetching` for conditional rendering that causes mount/unmount (use `isPending`)
- `staleTime` without matching `gcTime`
- Setting `staleTime` without explicit justification

**Exception**: `gcTime: 0` for polling queries to prevent cached terminal states.

## Key Patterns

### Feature Barrel Export

```tsx
// features/analysis/index.ts
export { StatsCard } from "./components/stats-card";
export { useAnalysis } from "./hooks/use-analysis";
```

### shadcn/ui

```bash
pnpm dlx shadcn@latest add [component]
```

### Error/Loading Boundaries

Each route has `error.tsx` and `loading.tsx` for error handling and Suspense.

### Loading Feedback for Server Fetching

Every server fetching state MUST have visible loading feedback. Silent data replacement is prohibited.

- **Background refetch** (cached data visible + `isFetching`): Wrap with `<RefreshOverlay isRefreshing={...}>` from `@/components/feedback`
- **Mutations**: Use optimistic update (`onMutate` + `onError` rollback) for instant UI response
- **Initial load** (`isLoading`): Use `<LoadingFallback>` or skeleton
- Anti-flash: `RefreshOverlay` has built-in 150ms delay — sub-150ms fetches won't show overlay

### Styling

- Tailwind CSS 4 + CSS Variables
- `cn()` for class merging: `cn("base", conditional && "extra")`
