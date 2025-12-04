# SpecVital Web

A web application that analyzes test files from public GitHub repositories and visualizes test status.

## Features

- **GitHub Repository Analysis**: Automatic detection and analysis of test files in public repositories
- **Multi-Framework Support**: Jest, Vitest, Playwright, Go test support
- **Test Dashboard**: Test statistics and tree view visualization
- **Filter & Search**: Framework-based filtering, test name search

## Quick Start

### Prerequisites

- Node.js 20+
- Go 1.24+
- pnpm
- [just](https://github.com/casey/just) (task runner)
- [air](https://github.com/air-verse/air) (Go hot reload)

### Installation

```bash
# Install dependencies
just deps
```

### Environment Setup

**Backend** (`src/backend/.env`):

```env
GITHUB_TOKEN=your_github_token  # Optional: increases rate limit
PORT=3000                       # Optional: default 3000
ALLOWED_ORIGINS=http://localhost:5173  # Optional: CORS origins
```

**Frontend** (`src/frontend/.env`):

```env
NEXT_PUBLIC_API_URL=http://localhost:3000
```

### Development

```bash
# Run backend (port 3000)
just run backend

# Run frontend (port 5173) - in another terminal
just run frontend
```

Open http://localhost:5173

### Build & Test

```bash
# Build all
just build

# Run all tests
just test

# Individual targets
just build backend
just test frontend
```

## Project Structure

```
src/
├── backend/                 # Go API Server
│   ├── cmd/server/         # Entry point
│   ├── analyzer/           # Analysis module
│   ├── github/             # GitHub API client
│   ├── health/             # Health check
│   └── common/             # Shared utilities
└── frontend/               # Next.js 15 App
    ├── app/                # Pages (App Router)
    ├── components/         # UI components
    └── lib/                # API client, utilities
```

## API Documentation

See [docs/api.md](docs/api.md) for detailed API documentation.

## Tech Stack

**Backend**:

- Go 1.24
- Chi Router
- [specvital/core](https://github.com/specvital/core) - Test parser

**Frontend**:

- Next.js 16 (App Router, Turbopack)
- React 19
- TailwindCSS 4
- shadcn/ui

## License

MIT
