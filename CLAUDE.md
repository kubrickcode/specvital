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

## Development Commands

### justfile

- `just deps` - Install dependencies (pnpm install)
- `just lint` - Run all linters (justfile, config, go)
- `just lint go` - Format Go code with gofmt
- `just test` - Run Go tests
- `just release` - Trigger release (merge to release branch)

### Single Test Execution

```bash
go test ./pkg/parser/strategies/jest/... -run TestJest
go test ./pkg/domain -v
```

## Architecture

### Core Structure

```
pkg/
├── domain/           # Domain models (Inventory, TestFile, TestSuite, Test)
└── parser/
    ├── scanner.go    # Entry point: Scan(), DetectTestFiles()
    ├── detector.go   # Test file detection
    ├── treesitter.go # Tree-sitter parser pooling/caching
    └── strategies/   # Framework-specific parsers (Strategy pattern)
        ├── registry.go       # Global registry
        ├── jest/             # Jest parser
        ├── vitest/           # Vitest parser
        ├── playwright/       # Playwright parser
        ├── gotesting/        # Go testing parser
        └── shared/jstest/    # Shared JS test logic
```

### Strategy Pattern

- Each framework implements `strategies.Strategy` interface
- Auto-registered via `strategies.Register()` in `init()`
- Blank import required: `_ "github.com/specvital/core/pkg/parser/strategies/jest"`

### Concurrency Model

- `scanner.go`: Parallel parsing with errgroup + semaphore
- `parser_pool.go`: Tree-sitter parser reuse via sync.Pool
- `treesitter.go`: Query compilation caching
