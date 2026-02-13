# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Infra - Database schema management and deployment configuration for the SpecVital monorepo

- **Purpose**: Database schema (Atlas HCL) + Railway deployment configs
- **No application code**: Only infrastructure definitions
- **Monorepo path**: `infra/` within the root monorepo

## Directory Structure

```
infra/
├── db/
│   ├── schema/
│   │   ├── schema.hcl          # Source of truth (declarative)
│   │   ├── migrations/         # Timestamped SQL diffs
│   │   └── rollbacks/
│   ├── atlas.hcl               # Environment configs (local/ci/prod)
│   └── docs/                   # Generated schema docs (tbls)
├── railway/                    # Centralized deployment configs
│   ├── web-backend/
│   ├── analyzer/
│   ├── spec-generator/
│   └── retention-cleanup/
└── justfile
```

## Commands

```bash
just --list              # View all available commands

# Migrations
just makemigration name  # Create new migration (Atlas diff → SQL)
just migrate             # Apply pending migrations
just reset               # Wipe DB + reapply all migrations

# Schema Visualization
just erd                 # Interactive ERD in browser (Atlas)
just gen-schema-docs     # Generate schema docs with tbls

# River Queue
just river-install       # Install River CLI
just river-dump          # Export River migration SQL
just river-migrate name  # Create River migration files (split for Atlas)
```

Note: Root `just migrate` orchestrates infra migration + schema dump for web/worker.

## Database Schema

**Tool**: [Atlas](https://atlasgo.io/) (HCL-based declarative schema)

### Schema Change Workflow

1. Edit `db/schema/schema.hcl`
2. `just makemigration <name>` - generates diff SQL
3. Review generated migration in `db/schema/migrations/`
4. `just migrate` - apply locally
5. Commit both schema.hcl and migration files

### Core Tables

- `codebases` - Repository metadata (host/owner/name)
- `analyses` - Test analysis runs with status tracking
- `test_suites` - Hierarchical test suites (self-referencing parent_id)
- `test_cases` - Individual tests with status/tags
- `users` / `oauth_accounts` - GitHub OAuth authentication

### Enums

- `analysis_status`: pending → running → completed/failed
- `test_status`: active, skipped, todo, focused, xfail
- `oauth_provider`: github

## Railway Deployment

All Railway deployment configs are centralized under `railway/`.

### Per-Service Files

| File           | Purpose                                              |
| -------------- | ---------------------------------------------------- |
| `railway.json` | Railway build/deploy config (dockerfilePath, region) |
| `Dockerfile`   | Multi-stage Go build (build context = repo root)     |

### Deployment Mechanism

CI copies `railway.json` to repo root before `railway up`:

```bash
cp infra/railway/<service>/railway.json railway.json
railway up --service <service> --detach
rm railway.json
```

`dockerfilePath` in `railway.json` points to `infra/railway/<service>/Dockerfile` (repo root relative).
Dockerfile `COPY` paths use `apps/web/...` or `apps/worker/...` prefixes (repo root as build context).

## Project-Specific Rules

### River Queue Integration

- River v0.26.0 for async job processing
- Migrations split into two files (schema + functions) for Atlas compatibility
- Use `just river-*` commands for River migration management

### Environment Variables

- `DATABASE_URL` - PostgreSQL connection (auto-configured in devcontainer)

### CI/CD

- `release.yml` - Semantic release on release branch
- `release-migrate.yml` - Production migration on release branch
- `release-deploy-frontend.yml` - Vercel frontend deployment
- `release-deploy-backend.yml` / `release-deploy-workers.yml` - Railway deployment
