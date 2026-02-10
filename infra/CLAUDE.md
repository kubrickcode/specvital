# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Infra - Shared infrastructure for the SpecVital platform (test file analysis for GitHub repositories)

- **Purpose**: Database schema management + documentation hub
- **No application code**: Only infrastructure and docs
- **Polyrepo foundation**: Must be started before collector/web repos

## Repository Map

```
specvital/
├── infra      # This repo - DB schema + docs (start first!)
├── collector  # Go Worker (consumes queue, uses core parser)
└── web        # NestJS + Next.js (API + Dashboard)
```

## Commands

```bash
just --list              # View all available commands

# Development
just deps                # Install dependencies (pnpm install)
just migrate             # Apply pending migrations
just reset               # Wipe DB + reapply all migrations
just makemigration name  # Create new migration

# Linting
just lint all            # Format justfile + config files
just lint config         # Prettier for JSON/YAML/MD
just lint justfile       # Format justfile

# Release
just release             # Trigger production release (main → release branch)

# River Queue (for collector integration)
just river-install       # Install River CLI
just river-dump          # Export River migration SQL
just river-migrate name  # Create River migration files

# Schema Visualization
just erd                 # Interactive ERD in browser (Atlas)
just install-tbls        # Install tbls (required for schema-doc)
just gen-schema-docs     # Generate schema docs with tbls (indexes, constraints, ERD)
```

## Database Schema

**Tool**: [Atlas](https://atlasgo.io/) (HCL-based declarative schema)

### Key Files

| Path                    | Description                               |
| ----------------------- | ----------------------------------------- |
| `db/schema/schema.hcl`  | Source of truth (declarative)             |
| `db/schema/migrations/` | Timestamped SQL migrations                |
| `db/atlas.hcl`          | Environment configs (local/ci/production) |

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

## Multi-Repo Development

### Network Setup

This repo creates `specvital-network`. Other repos connect to it:

```yaml
# In other repo's .devcontainer/docker-compose.yml
networks:
  specvital-network:
    name: specvital-network
    external: true
```

### Service Hostnames (from other devcontainers)

- PostgreSQL: `specvital-postgres:5432`
- Redis: `specvital-redis:6379`

## Documentation

| Path           | Content                       |
| -------------- | ----------------------------- |
| `docs/en/`     | English documentation         |
| `docs/ko/`     | Korean documentation          |
| `docs/en/adr/` | Architecture Decision Records |
| `docs/en/prd/` | Product Requirements          |

**Bilingual requirement**: Keep en/ko folders in sync when updating docs.

## Release Process

1. Merge to main via PR
2. `just release` merges main → release branch
3. GitHub Actions:
   - Analyzes commits (Conventional Commits)
   - Determines version bump (major/minor/patch)
   - Updates CHANGELOG.md
   - Applies migrations to production
   - Creates GitHub release
4. Release branch syncs back to main

## Railway Deployment

All Railway deployment configs are centralized under `infra/railway/`.

### Directory Structure

```
infra/railway/
├── web-backend/        # Go API server
├── analyzer/           # Test file analysis worker
├── spec-generator/     # AI spec document generation worker
└── retention-cleanup/  # Scheduled data cleanup (cron)
```

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

### Workflows

- `.github/workflows/release/deploy-web.yml` - web-backend deployment
- `.github/workflows/release/deploy-workers.yml` - worker services (matrix: analyzer, spec-generator, retention-cleanup)

## Project-Specific Rules

### River Queue Integration

- River v0.26.0 for async job processing
- Migrations split into two files (schema + functions) for Atlas compatibility
- Use `just river-*` commands for River migration management

### Environment Variables

- `DATABASE_URL` - PostgreSQL connection (auto-configured in devcontainer)
- `REDIS_URL` - Redis connection

### CI/CD

- `lint.yml` - Prettier check on PR/push
- `release.yml` - Semantic release on release branch
- `migrate.yml` - Production migration on release branch
