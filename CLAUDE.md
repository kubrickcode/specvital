# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Central documentation repository for the Specvital platform (https://specvital.github.io)

- VitePress-based static documentation site
- Multi-language support (English `/en/`, Korean `/ko/`)
- Auto-deployed to GitHub Pages on push to main branch

## Commands

```bash
# Dev server
pnpm docs:dev

# Build
pnpm docs:build

# Preview build
pnpm docs:preview

# Install dependencies
just deps

# Format all files
just lint all
```

## Architecture

```
docs/
├── .vitepress/config.mts    # VitePress config (navigation, sidebar)
├── en/                      # English docs
│   ├── prd/                 # Product Requirements Document
│   └── adr/                 # Architecture Decision Records
│       ├── core/            # Core microservice ADR
│       ├── collector/       # Collector microservice ADR
│       └── web/             # Web microservice ADR
└── ko/                      # Korean docs (mirrors en/ structure)
```

## Documentation Rules

### Adding New Documents

1. Create in both `/en/` and `/ko/` with identical structure
2. Update `docs/.vitepress/config.mts` (navigation/sidebar)
3. Update relevant `index.md` files in both languages

### Language Convention

- English docs (`/en/`): Write in English
- Korean docs (`/ko/`): Write in Korean
- Follow document language rules regardless of conversation language

### Korean Writing Style

- Use noun-ending style (명사형 종결): "~생성", "~필요", "~제공"
- Avoid verb-ending style (동사 종결어미): "~한다", "~이다"
- Example:
  - ✅ "정적 분석을 통한 테스트 카운트"
  - ❌ "정적 분석을 통해 테스트를 카운트한다"
