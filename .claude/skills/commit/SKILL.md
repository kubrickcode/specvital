---
name: commit
description: Generate Conventional Commits-compliant messages (feat/fix/docs/chore) in Korean and English. Use when you need to create a well-structured commit message for staged changes.
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git branch:*), Read, Grep, Write
disable-model-invocation: true
---

# Conventional Commits Message Generator

Generates commit messages following [Conventional Commits v1.0.0](https://www.conventionalcommits.org/) specification in both Korean and English. **Choose one version for your commit.**

## Repository State Analysis

- Git status: !`git status --porcelain`
- Current branch: !`git branch --show-current`
- Staged changes: !`git diff --cached --stat`
- Unstaged changes: !`git diff --stat`
- Recent commits: !`git log --oneline -10`

## What This Command Does

1. Checks current branch name to detect issue number (e.g., develop/shlee/32 → #32)
2. Checks which files are staged with git status
3. Performs a git diff to understand what changes will be committed
4. Generates commit messages in Conventional Commits format in both Korean and English
5. Adds "fix #N" at the end if branch name ends with a number
6. **Saves to commit_message.md file for easy copying**

## Conventional Commits Format (REQUIRED)

```
<type>[(optional scope)]: <description>

[optional body]

[optional footer: fix #N]
```

### Available Types

Analyze staged changes and suggest the most appropriate type:

| Type         | When to Use                                           | SemVer Impact |
| ------------ | ----------------------------------------------------- | ------------- |
| **feat**     | New feature or capability added                       | MINOR (0.x.0) |
| **fix**      | User-facing bug fix                                   | PATCH (0.0.x) |
| **ifix**     | Internal/infrastructure bug fix (CI, build, deploy)   | PATCH (0.0.x) |
| **perf**     | Performance improvements                              | PATCH         |
| **docs**     | Documentation only changes (README, comments, etc.)   | PATCH         |
| **style**    | Code formatting, missing semicolons (no logic change) | PATCH         |
| **refactor** | Code restructuring without changing behavior          | PATCH         |
| **test**     | Adding or fixing tests                                | PATCH         |
| **chore**    | Build config, dependencies, tooling updates           | PATCH         |
| **ci**       | CI/CD configuration changes                           | PATCH         |

**BREAKING CHANGE**: MUST use both type! format (exclamation mark after type) AND BREAKING CHANGE: footer with migration guide for major version bump.

### Confusing Cases: fix vs ifix vs chore

**Key distinction**: Does it affect **users** or only **developers/infrastructure**?

| Scenario                                              | Type       | Reason                                       |
| ----------------------------------------------------- | ---------- | -------------------------------------------- |
| Backend GitHub Actions test workflow not running      | `ifix`     | Bug in CI/CD that blocks development         |
| OOM error causing deployment failure                  | `ifix`     | Infrastructure bug blocking release          |
| E2E test flakiness causing false negatives            | `ifix`     | Testing infrastructure bug                   |
| Vite build timeout in production build                | `ifix`     | Build system bug                             |
| API returns 500 error for valid requests              | `fix`      | Users experience error responses             |
| Page loading speed improved from 3s to 0.8s           | `perf`     | Users directly feel the improvement          |
| App crashes when accessing profile page               | `fix`      | Users experience crash                       |
| Internal database query optimization (no user impact) | `refactor` | Code improvement, no measurable user benefit |
| Dependency security patch (CVE fix)                   | `chore`    | Build/tooling update (not a bug fix)         |
| Upgrading React version for new features              | `chore`    | Dependency update (not a bug fix)            |

## Commit Message Content Rules

**Commit messages are project history.** They answer "Why was this change worth making?" — never "What code was modified?"

### Subject Line: Purpose, Not Technical Action

The `<description>` MUST describe the **purpose or problem solved**, not what you technically did.

| Anti-pattern (technical action)                     | Correct (purpose/problem)                          |
| --------------------------------------------------- | -------------------------------------------------- |
| `getUser 함수에서 null 체크 누락 수정`              | `로그인 직후 프로필 페이지 접근 시 화면 깨짐 수정` |
| `상품 목록 API에 Redis 캐시 레이어 추가`            | `상품 목록 페이지 초기 로딩 속도 개선`             |
| `UserService를 UserRepository와 UserUseCase로 분리` | `사용자 도메인 로직과 DB 접근 책임 분리`           |
| `프로필 페이지에 비밀번호 변경 폼 컴포넌트 추가`    | `사용자가 직접 비밀번호를 변경할 수 있도록 지원`   |

**Litmus test**: If the subject ends with "추가", "수정", "변경", "개선" preceded by a technical noun — it's likely describing WHAT, not WHY. Rewrite to state the purpose.

### Body: Narrative Prose, Not Bullet Lists

The body is a **short narrative** (1-4 sentences) answering:

1. **왜 이 문제/필요가 존재했는가?** (맥락, 근본 원인)
2. **어떤 접근으로 해결했고, 왜 그 방식인가?** (접근 방식과 선택 근거를 한 흐름으로)

**FORBIDDEN**: Bullet-listed technical changes in body. That's what `git diff` provides.

**Exception**: Large-scope commits (5+ files, multiple concerns) may append a brief summary list AFTER the narrative body, prefixed with "주요 변경:".

### Korean Body Style Rules

**Sentence ending**: Use nominal endings (`~함`, `~됨`, `~필요함`, `~없음`, `~있음`), NOT declarative endings (`~한다`, `~된다`, `~했다`).

| Bad (declarative)                  | Good (nominal)                   |
| ---------------------------------- | -------------------------------- |
| 도메인 모델이 먼저 정의되어야 한다 | 도메인 모델이 먼저 정의되어야 함 |
| exact version으로 고정했다         | exact version으로 고정함         |
| 이 방식이 가장 적합하다고 판단했다 | 이 방식이 가장 적합하다고 판단함 |

**Line breaks**: Break ONLY after a sentence ends (after nominal ending like `~함.`, `~됨.`). NEVER break mid-sentence — if a clause continues, keep it on the same line.

| Bad (mid-sentence break)                                                     | Good (single line until sentence ends)                                  |
| ---------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| 후속 커밋이 공유 타입 위에 빌드되려면<br>도메인 모델이 먼저 정의되어야 한다. | 후속 커밋이 공유 타입 위에 빌드되려면 도메인 모델이 먼저 정의되어야 함. |

### Complexity-Based Formats

**Trivial** (typo, formatting, simple config):

```
type: brief description
```

**Simple**:

```
type: purpose/problem

Root cause or context in one sentence.
Approach taken and why (if non-obvious).
```

**Standard**:

```
type: purpose/problem

Why this problem existed (context).
How it was addressed and why this approach (rationale).

fix #N
```

**Large-scope** (exception — only case where list is allowed):

```
type: high-level purpose

Narrative context, approach, and rationale (2-4 sentences).

주요 변경:
- brief item 1
- brief item 2

fix #N
```

## Output Format

The command will provide:

1. Analysis of the staged changes (or all changes if nothing is staged)
2. **Creates commit_message.md file** containing both Korean and English versions
3. Copy your preferred version from the file

## Important Notes

- This command ONLY generates commit messages — it never performs actual commits
- **commit_message.md file contains both versions** — choose the one you prefer
- **Write the purpose, not the changelog** — subject = why it matters, body = context
- Branch issue numbers (e.g., develop/32) will automatically append "fix #N"
- Copy message from generated file and manually execute `git commit`

## Execution Instructions

### Phase 1: Synthesize Purpose (BEFORE analyzing diff)

1. Scan conversation history for the problem being solved or goal being pursued
2. Infer intent from branch name, changed file names, test descriptions
3. **Formulate a one-sentence purpose**: "이 커밋은 \_\_\_을/를 위해 필요하다"
   - This sentence becomes the seed for the subject line

### Phase 2: Analyze Changes

4. Run git commands to see staged changes (or all if none staged)
5. Verify: does the purpose from Phase 1 match actual changes? Adjust if needed

### Phase 3: Write Message

6. **Subject = purpose** from Phase 1, not a list of what was done
7. **Body = narrative** answering "why this problem?" + "why this solution?"
8. Choose template complexity based on change scope
9. **Select type** from the purpose and body content (not from file patterns)
10. Determine scope if needed (e.g., `fix(api):`, `feat(ui):`)

### Phase 4: Self-Verification (Hard Gate)

11. **Subject check**: Does it describe a technical action or a purpose?
    - Ends with "추가/수정/변경/개선" after a technical noun → REWRITE
12. **Body check**: Is it narrative prose or a bullet list of changes?
    - Contains `- ` bullet items listing technical changes → REWRITE as narrative
    - Exception: large-scope "주요 변경:" section after narrative
13. **History check**: "6개월 후 이 메시지만 읽고 왜 이 변경이 필요했는지 이해할 수 있는가?"
    - No → add context about the problem being solved

### Output

14. Generate both Korean and English versions
15. Write to `commit_message.md`
16. Present suggested type with reasoning
