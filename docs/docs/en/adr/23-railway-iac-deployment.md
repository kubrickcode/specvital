---
title: Railway IaC Deployment
description: ADR on Infrastructure as Code deployment using Railway CLI and GitHub Actions
---

# ADR-23: Railway IaC Deployment

> [한국어 버전](/ko/adr/23-railway-iac-deployment.md)

| Date       | Author       | Repos              |
| ---------- | ------------ | ------------------ |
| 2026-02-02 | @KubrickCode | worker, web, infra |

## Status

**Accepted** - Complements [ADR-06: PaaS-First Infrastructure](/en/adr/06-paas-first-infrastructure.md) and [ADR-22: Scheduler Removal and Railway Cron Migration](/en/adr/22-scheduler-removal-railway-cron.md)

## Context

### Problem Statement

Railway platform integration initially used the GitHub App automatic deployment feature, which presented several operational challenges:

1. **Manual Dashboard Configuration**: Each service required manual setup through Railway's web dashboard after GitHub App integration
2. **Non-Reproducible Deployments**: Service recreation required re-doing all manual configurations with no history tracking
3. **BuildArgs Limitation**: Railway's `railway.json` schema does not officially support `buildArgs` for Docker builds
4. **Environment Variable Rate Limiting**: Multiple sequential `railway variables --set` CLI calls triggered Railway API rate limits

### Why Change Was Required

| Issue                      | Impact                                | Frequency           |
| -------------------------- | ------------------------------------- | ------------------- |
| Manual Dashboard setup     | Deployment drift, onboarding friction | Every new service   |
| Non-reproducibility        | Disaster recovery risk                | Service recreation  |
| BuildArgs unsupported      | Cannot share Dockerfile               | Build configuration |
| Rate limiting on variables | Failed deployments, retry overhead    | Every deployment    |

## Decision

**Adopt GitHub Actions with Railway CLI for Infrastructure as Code (IaC) deployment, using per-service Dockerfiles and railway.json configuration files.**

### Configuration Structure

```
infra/
├── analyzer/
│   ├── Dockerfile          # Service-specific Dockerfile
│   └── railway.json        # Service configuration with $schema
├── spec-generator/
│   ├── Dockerfile
│   └── railway.json
├── retention-cleanup/
│   ├── Dockerfile
│   └── railway.json
└── web-backend/
    ├── Dockerfile
    └── railway.json
```

### Key Patterns

**1. Per-Service Dockerfile** (avoids buildArgs limitation):

Each service has its own Dockerfile instead of a shared Dockerfile with build arguments:

```dockerfile
# infra/analyzer/Dockerfile
FROM golang:1.22-alpine AS builder
# ... service-specific build steps
```

**2. railway.json Configuration**:

```json
{
  "$schema": "https://railway.com/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "infra/analyzer/Dockerfile"
  },
  "deploy": {
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 3
  }
}
```

**3. Cron Service Configuration** (for scheduled jobs):

```json
{
  "$schema": "https://railway.com/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "infra/retention-cleanup/Dockerfile"
  },
  "deploy": {
    "cronSchedule": "0 3 * * *",
    "restartPolicyType": "NEVER"
  }
}
```

**4. GitHub Actions Matrix Deployment**:

```yaml
deploy:
  strategy:
    matrix:
      service: [analyzer, spec-generator, retention-cleanup]
    fail-fast: true
  steps:
    - name: Install Railway CLI
      run: npm install -g @railway/cli
    - name: Deploy Service
      run: |
        cp infra/${{ matrix.service }}/railway.json railway.json
        railway up --service ${{ matrix.service }} --detach
        rm railway.json
      env:
        RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

## Options Considered

### Option A: GitHub Actions + Railway CLI (Selected)

Deploy via GitHub Actions workflow using Railway CLI with per-service configuration files.

**Pros:**

- Version controlled configuration in Git (code review, history)
- Reproducible deployments from any environment
- Per-service Dockerfile eliminates buildArgs limitation
- Matrix strategy enables parallel service deployments
- `$schema` provides IDE autocompletion and validation
- Consolidated API calls reduce rate limiting issues

**Cons:**

- Configuration duplication across service Dockerfiles
- Deployment tied to GitHub Actions availability
- Multiple railway.json files to maintain

### Option B: Railway GitHub App Integration (Previous)

Automatic deployment triggered by GitHub push events via Railway's native integration.

**Pros:**

- Zero CI/CD setup required
- Automatic deployment on push
- Simple mental model

**Cons:**

- Manual Dashboard configuration required per service
- Configuration not version controlled
- Environment drift between staging and production
- Same buildArgs limitation applies

### Option C: Pulumi/Terraform with Railway Provider

Use infrastructure-as-code tools with Railway's provider.

**Pros:**

- Full declarative infrastructure management
- State tracking and drift detection
- Cross-provider resource management

**Cons:**

- Additional tooling complexity
- State file management overhead
- Over-engineering for single-platform deployment
- Railway provider may have feature gaps

## Implementation

### Deployment Workflow

1. **Build Phase**: Railway CLI reads `railway.json` from repository root
2. **Config Copy**: Workflow copies service-specific config to root
3. **Deploy**: `railway up` executes with service name
4. **Cleanup**: Remove copied config file

### Environment Variables

- **Project-level**: Shared across services (DATABASE_URL, etc.)
- **Service-level**: Set via consolidated workflow step to avoid rate limiting

### Railway Project Token

Uses Project Token (not User Token) for CI/CD:

```yaml
env:
  RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

No `railway link` command needed with Project Token.

## Consequences

### Positive

**Reproducibility:**

- All deployment configuration version-controlled
- Any service can be recreated from repository alone
- Deployment changes visible in pull requests

**Developer Experience:**

- Schema validation provides IDE autocompletion
- Clear per-service configuration boundaries
- Familiar Git-based workflow

**Operational Reliability:**

- Consolidated API calls avoid rate limiting
- Matrix strategy enables parallel deployments
- Clear deployment logs in GitHub Actions

### Negative

**Configuration Duplication:**

- Per-service Dockerfiles contain similar base layers
- Minor changes may require updating multiple files
- **Mitigation**: Use multi-stage builds with common base image

**GitHub Actions Dependency:**

- Deployment requires GitHub Actions runners
- GitHub outages block automated deployments
- **Mitigation**: Railway CLI can be run manually as fallback

**Learning Curve:**

- Team must understand railway.json schema
- GitHub Actions workflow syntax required
- **Mitigation**: Centralized workflow with clear documentation

## Related ADRs

| ADR                                                                          | Relationship                           |
| ---------------------------------------------------------------------------- | -------------------------------------- |
| [ADR-06: PaaS-First Infrastructure](/en/adr/06-paas-first-infrastructure.md) | Foundation - Railway as PaaS platform  |
| [ADR-22: Scheduler Removal](/en/adr/22-scheduler-removal-railway-cron.md)    | Extension - Railway Cron configuration |

## References

- Commit `b842318`: Migrate to Railway IaC with GitHub Actions deployment (worker)
- Commit `c9ee36f`: Migrate Railway deployment to IaC (web)
- [Railway IaC Documentation](https://docs.railway.app/guides/config-as-code)
- [Railway CLI Reference](https://docs.railway.app/reference/cli-api)
