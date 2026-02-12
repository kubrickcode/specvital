# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Worker - Background job processing service for analyzing test files in GitHub repositories

- Queue-based async worker (River on PostgreSQL)
- Parser: `lib/` (monorepo shared library, referenced via Go `replace` directive)

### Workers

| Worker         | Kind                | Description                                    |
| -------------- | ------------------- | ---------------------------------------------- |
| AnalyzeWorker  | `analysis:analyze`  | Parse test files from GitHub repos             |
| SpecViewWorker | `specview:generate` | AI-powered test spec documentation (see below) |

### SpecView Worker

Generates human-readable spec documents from test files using Gemini AI.

- **Phase 1**: Domain/feature classification (gemini-2.5-flash)
- **Phase 2**: Test name → behavior conversion (gemini-2.5-flash-lite, parallel)
- **Cache**: Content hash-based deduplication
- **Reliability**: Circuit breaker, rate limiting, exponential backoff

Required env vars:

- `GEMINI_API_KEY`: Gemini API key
- `GEMINI_PHASE1_MODEL`: Phase 1 model (default: gemini-2.5-flash)
- `GEMINI_PHASE2_MODEL`: Phase 2 model (default: gemini-2.5-flash-lite)

## Documentation Map

| Context                         | Reference             |
| ------------------------------- | --------------------- |
| Architecture / Data flow        | `docs/en/`            |
| Design decisions (why this way) | `docs/en/adr/worker/` |
| Coding rules / Test patterns    | Root `.claude/rules/` |

## Commands

Before running commands, read `justfile` or check available commands via `just --list`

## Project-Specific Rules

### Auto-Generated Files (NEVER modify)

- `internal/infra/db/{queries.sql.go,models.go,db.go}`
- Workflow: `just dump-schema` → `just gen-sqlc`

### Monorepo Shared Library

- Parsing logic lives in `lib/` (monorepo root, referenced via Go `replace` directive)
- For parser changes, modify `lib/` directly

### Build Artifacts Cleanup

- `just build` outputs binaries to `bin/` directory
- After build verification, ALWAYS clean up: `rm -rf bin/`
- NEVER commit `bin/` directory (already in .gitignore)

## Common Workflows

### DB Schema Changes

1. Edit `infra/db/schema/schema.hcl` (monorepo root)
2. `cd infra && just makemigration <name>` → `just migrate`
3. `just dump-schema` → `just gen-sqlc`
4. Update `adapter/repository/` implementation

### Adding New Worker

1. Define worker in `adapter/queue/`
2. Register in `app/container.go`
3. Write tests
