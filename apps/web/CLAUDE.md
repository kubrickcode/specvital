# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Web - Full-stack application for analyzing and visualizing test files from GitHub repositories

- **OpenAPI-First**: Single spec (`openapi.yaml`) generates both Go and TypeScript types
- **Backend-for-Frontend**: Next.js as BFF, all business logic in Go backend
- **Feature-Based**: Domain modules in both `modules/` (backend) and `features/` (frontend)

## Documentation Map

| Context                 | Reference                      |
| ----------------------- | ------------------------------ |
| **Backend specifics**   | `src/backend/CLAUDE.md`        |
| **Frontend specifics**  | `src/frontend/CLAUDE.md`       |
| Coding rules / Patterns | `.claude/rules/`               |
| API specification       | `src/backend/api/openapi.yaml` |
| Available commands      | `just --list`                  |

## Commands

Run `just --list` to see all available commands. Key commands:

```bash
just deps                    # Install all dependencies
just run backend             # Go server on :8000 (hot reload via air)
just run frontend            # Next.js on :5173 (turbopack)
just gen-api                 # Generate Go/TS types from OpenAPI spec
just test                    # Run all tests
just build                   # Build all
```

## Critical Rules

### Type Generation Chain

```
openapi.yaml → just gen-api → { server.gen.go, generated-types.ts }
```

Always modify `openapi.yaml` first, then regenerate. Never edit generated files.

### External Dependency

- Test file parsing: `github.com/kubrickcode/specvital/packages/core`
- For parser changes → open issue in core repo first, NOT here

### E2E Test Maintenance

When implementing changes that affect user-facing behavior:

- **MUST** add or update E2E tests in `e2e/` directory
- UI changes (buttons, forms, navigation) → update corresponding `.spec.ts`
- New pages/features → create new E2E test file
- API-dependent tests → use mocked version (`*.mocked.spec.ts`)

## Data Flow

```
User → GitHub URL → /analyze/{owner}/{repo}
     → Backend → GitHub API → specvital/core parser → Response
```

## Environment Variables

### Required (Root `.env`)

| Variable                    | Purpose                                |
| --------------------------- | -------------------------------------- |
| `GITHUB_TOKEN`              | GitHub API rate limit (60 → 5000/hour) |
| `JWT_SECRET`                | Auth token signing (32+ chars)         |
| `ENCRYPTION_KEY`            | Sensitive data encryption (32+ chars)  |
| `GITHUB_APP_ID`             | GitHub App ID for org repo access      |
| `GITHUB_APP_SLUG`           | GitHub App slug name                   |
| `GITHUB_APP_PRIVATE_KEY`    | GitHub App private key (PEM format)    |
| `GITHUB_APP_WEBHOOK_SECRET` | Webhook signature verification         |

For stack-specific variables, see:

- Backend: `src/backend/CLAUDE.md`
- Frontend: `src/frontend/CLAUDE.md`

## GitHub App Integration

OAuth Apps require organization admin approval. GitHub App is used for organization repository access.

- **Architecture**: OAuth (login/personal repos) + GitHub App (organization repos)
- **Setup**: See `src/backend/.env.example` for configuration
- **Module**: `modules/github-app/` (backend)
