# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**CRITICAL**

- Always update CLAUDE.md and README.md When changing a feature that requires major work or essential changes to the content of the document. Ignore minor changes.
- Never create branches or make commits autonomously - always ask the user to do it manually
- ⚠️ MANDATORY SKILL LOADING - BEFORE editing files, READ relevant skills:
  - .ts → typescript
  - .tsx → typescript + react
  - .go → golang
  - .test.ts, .spec.ts → typescript-test + typescript
  - .test.go, \_test.go → go-test + golang
  - .graphql, resolvers, schema → graphql + typescript
  - package.json, go.mod → dependency-management
  - Path-based (add as needed): apps/web/** → nextjs, apps/api/** → nestjs
  - Skills path: .claude/skills/{name}/SKILL.md
  - 📚 REQUIRED: Display loaded skills at response END: `📚 Skills loaded: {skill1}, {skill2}, ...`
- If Claude repeats the same mistake, add an explicit ban to CLAUDE.md (Failure-Driven Documentation)
- Follow project language conventions for ALL generated content (comments, error messages, logs, test descriptions, docs)
  - Check existing codebase to detect project language (Korean/English/etc.)
  - Do NOT mix languages based on conversation language - always follow project convention
  - Example: English project → `describe("User authentication")`, NOT `describe("사용자 인증")`
- Respect workspace tooling conventions
  - Always use workspace's package manager (detect from lock files: pnpm-lock.yaml → pnpm, yarn.lock → yarn, package-lock.json → npm)
  - Prefer just commands when task exists in justfile or adding recurring tasks
  - Direct command execution acceptable for one-off operations
- Dependencies: exact versions only (`package@1.2.3`), forbid `^`, `~`, `latest`, ranges
  - New installs: check latest stable version first, then pin it (e.g., `pnpm add --save-exact package@1.2.3`)
  - CI must use frozen mode (`npm ci`, `pnpm install --frozen-lockfile`)
- Clean up background processes: always kill dev servers, watchers, etc. after use (prevent port conflicts)

**IMPORTANT**

- Next.js App Router: Client components as leaf nodes
  - Server components by default - render static content on server
  - Client components (`"use client"`) only for interactive parts (forms, state, event handlers)
  - Static text (titles, descriptions) → server component, interactive logic → client component
- Avoid unfounded assumptions - verify critical details
  - Don't guess file paths - use Glob/Grep to find them
  - Don't guess API contracts or function signatures - read the actual code
  - Reasonable inference based on patterns is OK
  - When truly uncertain about important decisions, ask the user
- Always gather context before starting work
  - Read related files first (don't work blind)
  - Check existing patterns in codebase
  - Review project conventions (naming, structure, etc.)
- Always assess issue size and scope accurately - avoid over-engineering simple tasks
  - Apply to both implementation and documentation
  - Verbose documentation causes review burden for humans

## Commands

```bash
# Dependencies
just deps                    # Install all dependencies

# Development
just run backend             # Go server on :8000 (uses air for hot reload)
just run frontend            # Next.js on :5173 (uses turbopack)

# Build & Test
just build                   # Build all
just build backend           # go build ./...
just build frontend          # pnpm build
just test                    # Test all
just test backend            # go test -v ./...
just test frontend           # pnpm test (vitest)

# Code Generation
just gen-api                 # Generate Go/TS types from OpenAPI spec

# Lint
just lint                    # All targets
just lint go                 # gofmt
```

## Architecture

### Backend (Go + Chi + OpenAPI)

**Architecture**: OpenAPI-First with Strict Server Interface

```
src/backend/
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── api/types.gen.go    # Generated from OpenAPI spec (DO NOT EDIT)
│   └── app/app.go          # Router setup with strict handler
├── api/
│   ├── openapi.yaml        # API spec (single source of truth)
│   └── oapi-codegen.yaml   # Code generation config
├── modules/
│   └── analyzer/           # Core analysis module
│       ├── server.go       # StrictServerInterface implementation
│       └── service.go      # Business logic (GitHub → Core parser)
├── github/                 # GitHub API client
│   ├── client.go           # ListFiles, GetFileContent, rate limit tracking
│   └── errors.go           # ErrNotFound, ErrForbidden, ErrRateLimited
└── common/
    ├── dto/problem.go      # RFC 7807 error responses
    └── middleware/         # CORS, Logger, Compress
```

**Type Generation**: OpenAPI spec → Go types (compile-time verified) + TS types

### Frontend (Next.js 16 App Router)

```
src/frontend/
├── app/[locale]/                       # i18n routing (next-intl)
│   ├── page.tsx                        # Home - URL input form
│   ├── layout.tsx                      # Locale layout with Header
│   └── analyze/[owner]/[repo]/
│       ├── page.tsx                    # Dashboard - test results
│       ├── error.tsx                   # Error boundary
│       └── loading.tsx                 # Loading skeleton
│
├── features/                           # Domain modules
│   ├── analysis/                       # Analysis feature
│   │   ├── components/                 # stats-card, test-list, etc.
│   │   └── index.ts                    # Barrel export
│   └── home/                           # Home feature
│       ├── components/                 # url-input-form
│       ├── lib/                        # github-url validation
│       └── index.ts
│
├── components/                         # Shared components
│   ├── ui/                             # shadcn/ui primitives
│   ├── layout/                         # header
│   ├── theme/                          # theme-provider, toggle, language-selector
│   ├── feedback/                       # rate-limit-warning
│   └── skeletons/                      # Loading skeletons
│
├── lib/                                # Shared utilities
│   ├── api/                            # API client
│   │   ├── generated-types.ts          # Generated from OpenAPI (DO NOT EDIT)
│   │   ├── types.ts                    # Type re-exports for convenience
│   │   └── client.ts                   # API client implementation
│   ├── utils/                          # cn, formatAnalysisDate
│   └── styles/                         # framework-colors
│
├── i18n/                               # Internationalization config
└── messages/                           # i18n messages (en.json, ko.json)
```

#### Adding New Features

```typescript
// 1. Create feature folder
features/[feature-name]/
├── components/
├── api/           # (optional)
├── hooks/         # (optional)
└── index.ts       # Barrel export

// 2. Import from feature
import { Component } from "@/features/[feature-name]";
```

### Data Flow

1. User enters GitHub URL → validate → navigate to `/analyze/{owner}/{repo}`
2. Frontend calls `GET /api/analyze/{owner}/{repo}`
3. Backend fetches repo tree via GitHub API → filters test files → parses with specvital/core
4. Returns AnalysisResult with TestSuites and Summary

## Environment Variables

**Backend** (`src/backend/.env`):

- `GITHUB_TOKEN` - increases rate limit (60 → 5000/hour)
- `PORT` - default 8000
- `ALLOWED_ORIGINS` - CORS origins

**Frontend** (`src/frontend/.env`):

- `NEXT_PUBLIC_API_URL` - backend URL (default: http://localhost:8000)
