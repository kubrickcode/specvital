# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SpecVital Docs - VitePress documentation site for the SpecVital platform

- **Monorepo path**: `docs/` within the root monorepo
- VitePress-based static site, deployed to GitHub Pages
- Multi-language support (English `/en/`, Korean `/ko/`)

## Commands

```bash
# Dev server
pnpm docs:dev

# Build
pnpm docs:build

# Preview build
pnpm docs:preview
```

## Architecture

```
docs/
├── .vitepress/config.mts    # VitePress config (navigation, sidebar)
├── en/                      # English docs
│   ├── prd/                 # Product Requirements Document
│   └── adr/                 # Architecture Decision Records
│       ├── XX-*.md          # Cross-cutting ADR (root level)
│       ├── core/            # Core (lib) ADR
│       ├── worker/          # Worker ADR
│       └── web/             # Web ADR
├── ko/                      # Korean docs (mirrors en/ structure)
└── schema/                  # OpenAPI schema definitions
```

## Documentation Rules

### CRITICAL: Infra ADR Placement

> **Infra has NO standalone ADR folder.**

`infra/` contains only shared database schemas and deployment configs. It has no independent business logic, so:

- **Database schema ADR** → Place in **Cross-cutting ADR** (root level `adr/XX-*.md`)
- **NEVER create** `adr/infra/` folder

**Rationale**: Infra serves all services (lib, worker, web). Schema decisions are cross-cutting by nature.

### CRITICAL: Adding New Documents

> **VIOLATION WARNING: Skipping ANY file is a critical failure.**
>
> **`docs/en/index.md` and `docs/ko/index.md` are ALWAYS required but MOST MISSED (10+ occurrences).**

**Step 1: Create Document Files**

- Create in both `/en/` and `/ko/` with identical structure
- Follow naming convention (e.g., `01-xxx.md`, `02-yyy.md`)

**Step 2: Update Sidebar Navigation**

- Update `.vitepress/config.mts` - add to sidebar configuration

**Step 3: Update Index Files (COMPLETE CHECKLIST)**

> **CRITICAL**: Use this checklist for EVERY new document. Missing ANY file causes broken links.

#### For ADR Documents (e.g., `adr/core/13-xxx.md`)

**ALWAYS update these 6 files** (check each box):

```
☐ docs/en/index.md                 ← Homepage (ADR links section)
☐ docs/ko/index.md                 ← Homepage (ADR links section)
☐ docs/en/adr/index.md             ← ADR overview (repository section)
☐ docs/ko/adr/index.md             ← ADR overview (repository section)
☐ docs/en/adr/{category}/index.md  ← Category index (core/worker/web)
☐ docs/ko/adr/{category}/index.md  ← Category index (core/worker/web)
```

**Concrete examples:**

- For `adr/core/13-xxx.md`: Update 6 files (both homepage + adr/index + adr/core/index)
- For `adr/worker/08-xxx.md`: Update 6 files (both homepage + adr/index + adr/worker/index)
- For `adr/web/02-xxx.md`: Update 6 files (both homepage + adr/index + adr/web/index)

#### For PRD Documents (e.g., `prd/07-xxx.md`)

**ALWAYS update these 4 files** (check each box):

```
☐ docs/en/index.md       ← Homepage (PRD links section)
☐ docs/ko/index.md       ← Homepage (PRD links section)
☐ docs/en/prd/index.md   ← PRD index
☐ docs/ko/prd/index.md   ← PRD index
```

**Step 4: Verification (File Count Check)**

| Document Type         | Files to Update | Breakdown                                   |
| --------------------- | --------------- | ------------------------------------------- |
| ADR (core/worker/web) | **6 files**     | 2 homepage + 2 adr/index + 2 category/index |
| Root ADR (adr/XX.md)  | **4 files**     | 2 homepage + 2 adr/index                    |
| PRD                   | **4 files**     | 2 homepage + 2 prd/index                    |

> **RED FLAG CHECK:**
>
> - Updated only 4 files for category ADR → **MISSED 2 FILES** (homepage)
> - Didn't touch `docs/en/index.md` → **CRITICAL MISS**
> - Didn't touch `docs/ko/index.md` → **CRITICAL MISS**
>
> **If unsure, ALWAYS update `docs/en/index.md` and `docs/ko/index.md` FIRST.**

### Language Convention

- English docs (`/en/`): Write in English
- Korean docs (`/ko/`): Write in Korean
- Follow document language rules regardless of conversation language

### CRITICAL: Korean Writing Style

> **MUST use noun-ending style. Verb endings are prohibited.**

- Use noun-ending style (명사형 종결): "~생성", "~필요", "~제공", "~처리", "~확인"
- **NEVER use verb-ending style** (동사 종결어미): "~한다", "~이다", "~합니다", "~됩니다"
- Example:
  - "정적 분석을 통한 테스트 카운트"
  - "데이터 무결성 검증"
  - "API 호출 최소화"
