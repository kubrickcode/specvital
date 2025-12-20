# Specvital Infrastructure

Shared infrastructure for Specvital project - local development environment and database schema management.

## Quick Start

### VS Code Devcontainer (Recommended)

1. Open this folder in VS Code
2. Click "Reopen in Container" when prompted
3. PostgreSQL and Redis start automatically

### Docker Compose

```bash
cd .devcontainer && docker compose up -d
```

## Services

| Service    | Port | Container Name     |
| ---------- | ---- | ------------------ |
| PostgreSQL | 5432 | specvital-postgres |
| Redis      | 6379 | specvital-redis    |

## Commands

```bash
just --list              # View all commands

just deps                # Install dependencies
just migrate             # Apply DB migrations
just reset               # Reset DB (wipe + migrate)
just makemigration name  # Create new migration
just lint all            # Format all files
```

## Multi-Repository Setup

```
specvital/
├── infra        # This repo - start first!
├── collector    # Go Worker
└── web          # NestJS + Next.js
```

**Workflow**: Open infra devcontainer first → then open other repos

### Connecting from Other Repos

Add to `.devcontainer/docker-compose.yml`:

```yaml
networks:
  specvital-network:
    name: specvital-network
    external: true
```

Environment variables:

```
DATABASE_URL=postgres://postgres:postgres@specvital-postgres:5432/specvital?sslmode=disable
REDIS_URL=redis://specvital-redis:6379
```

## Schema Management

Database schema is managed with [Atlas](https://atlasgo.io/).

| Path                   | Description           |
| ---------------------- | --------------------- |
| `db/schema/schema.hcl` | Declarative schema    |
| `db/schema/migrations` | Timestamped SQL files |
| `db/atlas.hcl`         | Environment configs   |

## Troubleshooting

**Network not found**: Start infra devcontainer first, or run `docker network create specvital-network`

**Port conflict**: `docker ps | grep -E "5432|6379"` → stop conflicting container
