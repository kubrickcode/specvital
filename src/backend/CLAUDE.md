# Backend CLAUDE.md

Go API server with Chi router and OpenAPI-first architecture.

## Architecture: 5-Layer Clean Architecture

```
modules/{module}/
├── domain/
│   ├── entity/      # Business entities (pure Go)
│   └── port/        # Interface definitions (DIP)
├── usecase/         # Business logic (1 file per feature)
├── adapter/         # External implementations (DB, API, Queue)
│   └── mapper/      # Domain ↔ API conversion
└── handler/         # HTTP entry points
```

### Dependency Direction

```
handler → usecase → domain ← adapter
                      ↑
              (implements port)
```

### Layer Responsibilities

| Layer             | Responsibility                                   | Import Rules                |
| ----------------- | ------------------------------------------------ | --------------------------- |
| **domain/entity** | Business entities, value objects (pure Go)       | No external dependencies    |
| **domain/port**   | Interface definitions (Repository, Client, etc.) | Only entity                 |
| **usecase**       | Business logic, orchestration                    | Only domain (entity + port) |
| **adapter**       | External implementations (DB, API, Queue)        | domain + external libraries |
| **handler**       | HTTP entry points, request/response handling     | usecase + adapter/mapper    |

### Error Handling Pattern

```go
// Usecase returns domain error
if analysis == nil {
    return nil, domain.ErrNotFound
}

// Handler maps to HTTP status
switch {
case errors.Is(err, domain.ErrNotFound):
    return api.GetAnalysis404JSONResponse{...}, nil
}
```

## Auto-Generated Files (NEVER modify)

| File                         | Source             | Regenerate      |
| ---------------------------- | ------------------ | --------------- |
| `internal/api/server.gen.go` | `api/openapi.yaml` | `just gen-api`  |
| `internal/db/*.sql.go`       | `queries/*.sql`    | `just gen-sqlc` |

## Common Workflows

### API Change

1. Modify `api/openapi.yaml`
2. Run `just gen-api`
3. Implement handler method (satisfies `StrictServerInterface`)
4. Add usecase logic with domain types
5. Create mapper if new response type

### Database Change

1. Add/modify SQL in `queries/*.sql`
2. Run `just gen-sqlc`
3. Update adapter/repository to use new generated methods

### New Module

Create `modules/[name]/` with Clean Architecture structure:

```
modules/[name]/
├── domain/
│   ├── entity/
│   │   └── models.go      # Business entities
│   ├── port/
│   │   └── repository.go  # Interface definitions
│   └── errors.go          # Sentinel errors (ErrNotFound, etc.)
├── usecase/
│   └── {feature}.go       # One file per use case
├── adapter/
│   ├── repository_postgres.go  # Port implementation
│   └── mapper/
│       └── response.go    # Domain ↔ API conversion
└── handler/
    └── http.go            # StrictServerInterface implementation
```

### Adding a New Use Case

1. Define Input/Output types in usecase file
2. Implement usecase struct with port dependencies
3. Add constructor function
4. Wire in `common/server/app.go`

## Authentication

Hybrid Token Architecture (Access + Refresh Token)

```
GitHub OAuth → Access Token (15min JWT) + Refresh Token (7d, DB-stored)
                     ↓                           ↓
              httpOnly cookie              httpOnly cookie (Strict)
              (Lax, stateless)             (stateful, rotation)
```

### Token Flow

1. **Login**: OAuth callback issues both tokens, refresh token stored in DB
2. **API Request**: Middleware validates access token (stateless JWT check)
3. **Token Expiry**: Frontend auto-refreshes via `/api/auth/refresh` on 401
4. **Logout**: Both tokens cleared, refresh token revoked in DB

### Security Features

- **Token Rotation**: New refresh token on each use
- **Reuse Detection**: Family-based revocation on token replay attack
- **Hash Storage**: Only refresh token hash stored in DB

### Key Files

- Token generation: `modules/auth/adapter/token_manager_jwt.go`
- Refresh logic: `modules/auth/usecase/refresh_token.go`
- Cookie handling: `modules/auth/handler/http.go`
- Middleware: `common/middleware/auth.go`
- OAuth config: `internal/infra/oauth.go`

## Async Processing

- **Queue**: River (PostgreSQL-backed)
- **Workers**: `modules/analyzer/adapter/queue_river.go`
- Used for: Repository analysis tasks

## Environment Variables

| Variable                    | Purpose                      | Default |
| --------------------------- | ---------------------------- | ------- |
| `PORT`                      | Server port                  | 8000    |
| `ALLOWED_ORIGINS`           | CORS origins                 | -       |
| `DATABASE_URL`              | PostgreSQL connection        | -       |
| `GITHUB_TOKEN`              | GitHub API rate limit        | -       |
| `JWT_SECRET`                | Token signing (32+ chars)    | -       |
| `ENCRYPTION_KEY`            | Sensitive data (32+ chars)   | -       |
| `FRONTEND_URL`              | OAuth redirect               | -       |
| `GITHUB_APP_ID`             | GitHub App ID                | -       |
| `GITHUB_APP_SLUG`           | GitHub App slug name         | -       |
| `GITHUB_APP_PRIVATE_KEY`    | GitHub App private key (PEM) | -       |
| `GITHUB_APP_WEBHOOK_SECRET` | Webhook signature secret     | -       |

## Key Patterns

### StrictServerInterface

All handlers implement the generated interface:

```go
type StrictServerInterface interface {
    GetAnalysis(ctx context.Context, request GetAnalysisRequestObject) (GetAnalysisResponseObject, error)
    // ...
}
```

### Domain Errors → HTTP Status

```go
var (
    ErrNotFound      = errors.New("not found")       // → 404
    ErrAlreadyQueued = errors.New("already queued")  // → 409
    ErrRateLimited   = errors.New("rate limited")    // → 429
)
```

### Port Interface Pattern

```go
// domain/port/repository.go
type Repository interface {
    GetAnalysis(ctx context.Context, id string) (*entity.Analysis, error)
    SaveAnalysis(ctx context.Context, analysis *entity.Analysis) error
}

// adapter/repository_postgres.go
type PostgresRepository struct {
    queries *db.Queries
}

var _ port.Repository = (*PostgresRepository)(nil)  // Compile-time check
```

### UseCase Pattern

```go
// usecase/get_analysis.go
type GetAnalysisInput struct {
    Owner string
    Repo  string
}

type GetAnalysisOutput struct {
    Analysis   *entity.Analysis
    TestSuites []entity.TestSuite
}

type GetAnalysisUseCase struct {
    repo port.Repository
}

func (uc *GetAnalysisUseCase) Execute(ctx context.Context, input GetAnalysisInput) (*GetAnalysisOutput, error) {
    // Business logic here
}
```

### SQLc Usage

```go
// queries/analysis.sql
-- name: GetAnalysis :one
SELECT * FROM analyses WHERE id = $1;

// Generated: internal/db/analysis.sql.go
func (q *Queries) GetAnalysis(ctx context.Context, id int64) (Analysis, error)
```

## Import Rules (Enforced by depguard)

```
domain/entity  → No external dependencies
domain/port    → Only entity
usecase        → Only domain (entity + port)
adapter        → domain + external libraries
handler        → usecase + adapter/mapper
```

Cross-module dependencies use port interfaces (e.g., `auth.TokenProvider` used by `analyzer`).
