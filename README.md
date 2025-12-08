# Specvital Infrastructure

Shared infrastructure for Specvital project - provides local development environment and database schema management.

## Quick Start

### Option 1: VS Code Devcontainer (Recommended)

1. Open this folder in VS Code
2. Click "Reopen in Container" when prompted
3. PostgreSQL and Redis start automatically

### Option 2: Docker Compose Directly

```bash
cd .devcontainer
docker compose up -d

# Check status
docker compose ps
```

## Services

| Service    | Port | Container Name      | Description   |
| ---------- | ---- | ------------------- | ------------- |
| PostgreSQL | 5432 | specvital-postgres  | Database      |
| Redis      | 6379 | specvital-redis     | Queue & Cache |

## Connection

### Inside Devcontainer

Environment variables are pre-configured:

```bash
# PostgreSQL
psql $DATABASE_URL

# Or directly
psql -h postgres -U postgres -d specvital
```

```bash
# Redis
redis-cli -h redis -p 6379
```

### From Host Machine

```bash
# PostgreSQL
psql -h localhost -U postgres -d specvital
# Password: postgres

# Redis
redis-cli -h localhost -p 6379
```

### From Other Devcontainers

Other Specvital repositories can connect via the shared `specvital-network`:

```bash
# PostgreSQL (use container name)
psql -h specvital-postgres -U postgres -d specvital

# Redis (use container name)
redis-cli -h specvital-redis -p 6379
```

## Environment Variables

### Inside Devcontainer (Auto-configured)

```
DATABASE_URL=postgres://postgres:postgres@postgres:5432/specvital?sslmode=disable
REDIS_URL=redis://redis:6379
```

### For Other Repos (via specvital-network)

```
DATABASE_URL=postgres://postgres:postgres@specvital-postgres:5432/specvital?sslmode=disable
REDIS_URL=redis://specvital-redis:6379
```

## Multi-Repository Setup

This infrastructure is shared across multiple Specvital repositories:

```
specvital/
├── infra        # This repo - start first!
├── collector    # Go Worker
└── web          # NestJS + Next.js
```

### Workflow

1. **Open infra in devcontainer first** (or run `docker compose up -d` in `.devcontainer/`)
2. **Open other repos in devcontainers** - they connect to the shared network

### Devcontainer Configuration for Consumer Repos

Add to `.devcontainer/docker-compose.yml`:

```yaml
services:
  workspace:
    # ... your config
    networks:
      - specvital-network

networks:
  specvital-network:
    name: specvital-network
    external: true
```

Or for simple devcontainer.json (without docker-compose):

```jsonc
{
  "name": "specvital-<repo>",
  "runArgs": ["--network=specvital-network"],
  "containerEnv": {
    "DATABASE_URL": "postgres://postgres:postgres@specvital-postgres:5432/specvital?sslmode=disable",
    "REDIS_URL": "redis://specvital-redis:6379"
  }
}
```

## Commands

```bash
# From .devcontainer directory
cd .devcontainer

# Start all services
docker compose up -d

# Stop all services
docker compose down

# Stop and remove volumes (reset data)
docker compose down -v

# View logs
docker compose logs -f postgres
docker compose logs -f redis

# Connect to PostgreSQL
docker compose exec postgres psql -U postgres -d specvital

# Connect to Redis
docker compose exec redis redis-cli
```

## Troubleshooting

### Network not found error

If other repos see `network specvital-network not found`:

```bash
# Make sure infra devcontainer is running first
# Or manually create the network:
docker network create specvital-network
```

### Port already in use

```bash
# Find conflicting containers
docker ps | grep -E "5432|6379"
docker stop <container_id>
```

### Connection refused from other devcontainer

1. Verify infra is running: `docker ps | grep specvital`
2. Check network exists: `docker network ls | grep specvital`
3. Verify your container is on the network: `docker network inspect specvital-network`

## Schema Management

Database schema is managed with [Atlas](https://atlasgo.io/). See `schema/` directory for details.

```bash
# Apply migrations
atlas migrate apply --url "$DATABASE_URL"
```
