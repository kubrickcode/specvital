# SpecVital

> Automated Test Inventory Generation via Static Code Analysis

**[specvital.com](https://specvital.com)** | [한국어](README.ko.md)

## What is SpecVital?

SpecVital is a platform that automatically generates test inventories from GitHub repositories using AST-based static analysis. Simply enter a repository URL and instantly see the complete test landscape—no CI/CD integration required.

### Core Value Proposition

| Feature            | Description                                        |
| ------------------ | -------------------------------------------------- |
| **Instant**        | Results in seconds, not hours of setup             |
| **Accurate**       | Tree-sitter AST parsing, deterministic (not AI)    |
| **Multi-Language** | Support for major test frameworks across languages |

## How It Works

```
GitHub URL → AST Parsing (Tree-sitter) → Test Inventory
```

1. Enter a public GitHub repository URL
2. SpecVital clones and parses test files using [specvital/core](https://github.com/specvital/core)
3. View test suites, cases, and structure in a visual dashboard

**Note**: This is static analysis—we parse test definitions, not execution results. You'll see what tests exist and their structure, not pass/fail status.

## Supported Frameworks

20+ test frameworks across multiple languages:

| Language              | Frameworks                               |
| --------------------- | ---------------------------------------- |
| JavaScript/TypeScript | Jest, Vitest, Playwright, Cypress, Mocha |
| Go                    | testing                                  |
| Python                | pytest, unittest                         |
| Java                  | JUnit 5, TestNG                          |
| Kotlin                | Kotest                                   |
| C#                    | NUnit, xUnit, MSTest                     |
| Ruby                  | RSpec, Minitest                          |
| PHP                   | PHPUnit                                  |
| Rust                  | cargo test                               |
| C++                   | Google Test                              |
| Swift                 | XCTest                                   |

See [specvital/core](https://github.com/specvital/core) for details.

## Current Status

SpecVital is in active development. Currently available:

- ✅ GitHub OAuth authentication
- ✅ Repository analysis via URL input
- ✅ Test tree visualization with statistics
- ✅ Multi-framework detection
- ✅ Filter and search capabilities
- ✅ Internationalization (English, Korean)

## Architecture

This repository contains the web application (Frontend + Backend API). The full system consists of:

| Repository                                                    | Role                                |
| ------------------------------------------------------------- | ----------------------------------- |
| [specvital/web](https://github.com/specvital/web)             | Web dashboard + REST API            |
| [specvital/core](https://github.com/specvital/core)           | Parser library (Tree-sitter based)  |
| [specvital/collector](https://github.com/specvital/collector) | Background worker for analysis jobs |
| [specvital/infra](https://github.com/specvital/infra)         | Database schema and infrastructure  |

## Use Cases

- **Engineering Managers**: Understand team's test coverage landscape
- **QA Leads**: Manage automated test inventory across projects
- **Staff Engineers**: Map existing tests before legacy refactoring

## Feedback

- Questions & Ideas: [GitHub Discussions](https://github.com/orgs/specvital/discussions)
- Bug Reports: [GitHub Issues](https://github.com/specvital/web/issues)

## License

MIT
