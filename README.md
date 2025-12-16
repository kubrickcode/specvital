# specvital collector

Repository analysis worker service for SpecVital platform.

## Architecture

Clean Architecture with separate entry points for Railway deployment:

```
src/cmd/
├── worker/      # Asynq worker - queue processing (Railway service #1)
├── scheduler/   # Cron scheduler - periodic jobs (Railway service #2)
├── enqueue/     # CLI tool for manual task enqueue
└── collector/   # Legacy entry point (deprecated, use worker)
```

## Build

```bash
# Build all binaries
just build

# Build specific target
just build worker
just build scheduler
just build enqueue

# Output: bin/worker, bin/scheduler, bin/enqueue, bin/collector
```

## Development

```bash
# Run worker locally with hot reload
just run local

# Run scheduler locally
just run-scheduler local

# Run tests
just test unit
just test integration
just test all
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string for Asynq

## Railway Deployment

Deploy as two separate services:

- **Worker**: `bin/worker` - processes analysis tasks from queue (scalable)
- **Scheduler**: `bin/scheduler` - runs periodic cron jobs (single instance)
