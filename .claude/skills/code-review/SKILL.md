---
name: code-review
description: Review current git changes or latest commit using code-reviewer and architect-reviewer agents. Use after completing code changes to get comprehensive quality feedback.
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git show:*), Task
disable-model-invocation: true
---

# Code Review Command

## Purpose

Perform comprehensive code review of current changes or latest commit. Automatically invokes:

- **code-reviewer agent**: Always (code quality, security, maintainability)
- **architect-reviewer agent**: Conditionally (architectural impact assessment)

---

## Workflow

### 1. Determine Review Target

**Check for uncommitted changes:**

```bash
git status --porcelain
```

**Decision:**

- If changes exist: Review uncommitted changes via `git diff HEAD` (staged + unstaged)
- If no changes: Review latest commit via `git log -1` and `git show HEAD`
- If user specified a commit hash: Review that commit via `git show <hash>`

**Empty diff guard:** If no diff output is produced, report "No changes to review" and stop.

### 2. Exclusion Patterns

Filter out the following from diff analysis and agent prompts:

- Lock files: `package-lock.json`, `pnpm-lock.yaml`, `yarn.lock`, `go.sum`, `Cargo.lock`
- Generated/build: `dist/`, `build/`, `*.generated.*`, `*.min.*`
- Source maps: `*.map`
- Binary/media: images, fonts, compiled artifacts

Use `git diff HEAD -- . ':!package-lock.json' ':!pnpm-lock.yaml' ':!yarn.lock' ':!go.sum' ':!Cargo.lock' ':!dist' ':!build'` to exclude these.

### 3. Analyze Change Scope

**Extract metrics from git output:**

Count:

- Files changed
- Lines added/removed
- Directories affected
- New files vs modifications

Detect:

- Database migration files (`migrations/`, `schema.prisma`, `*.sql`)
- API contract files (`*.graphql`, `openapi.yaml`, `*.proto`, API route files)
- Configuration files (`package.json`, `go.mod`, `pyproject.toml`, `Cargo.toml`, `build.gradle`, `docker-compose.yml`)
- Infrastructure files (`*.tf`, `Dockerfile`, K8s manifests, CI/CD configs)
- Security-sensitive files (`auth/`, `**/middleware/auth*`, `*.pem`, `*.key`)
- New directories or modules

### 4. Invoke code-reviewer Agent

**Always execute using Task tool with subagent_type: "code-reviewer":**

Provide structured context to avoid redundant discovery:

```
## Review Context

- **Target**: {uncommitted changes / commit <hash>}
- **Work Type**: {feature/bugfix/refactor/chore/prototype}
- **Commit Message**: {message from git log or commit_message.md}

## Scope Boundary

Architecture concerns (component boundaries, module coupling, system scalability, deployment architecture) are out of scope. Do NOT review these.

## Diff

{see diff size rules below}
```

**Diff size rules:**

- **≤ 500 lines**: Include full diff inline
- **> 500 lines**: Include `git diff --stat` summary only, and instruct: "Read specific files using Read tool as needed for detailed review."

### 5. Conditional architect-reviewer Invocation

**Scope Classification Criteria** (score each criterion met):

| Criterion             | Threshold                                             | Score |
| --------------------- | ----------------------------------------------------- | ----- |
| Files Modified        | ≥ 8 files                                             | +1    |
| Total Lines Changed   | ≥ 300 lines                                           | +1    |
| Directories Affected  | ≥ 3 directories                                       | +1    |
| New Module/Package    | New directory with ≥3 files                           | +1    |
| API Contract Changes  | Modified: schema/, .graphql, .proto, api/, routes     | +2    |
| Database Schema       | Modified: migrations, schema files, ORM models        | +2    |
| Config/Infrastructure | Modified: docker-compose, Dockerfile, K8s, CI/CD, .tf | +1    |
| Dependency Changes    | Modified: package.json, go.mod, requirements.txt      | +1    |
| Security-Sensitive    | Modified: auth/, middleware/auth*, *.pem, \*.key      | +2    |

**Decision**:

- **Score < 3** → Invoke `code-reviewer` only
- **Score ≥ 3** → Invoke both `code-reviewer` and `architect-reviewer` **in parallel** (two Task calls in one message)

**architect-reviewer structured context:**

```
## Review Context

- **Target**: {uncommitted changes / commit <hash>}
- **Work Type**: {feature/bugfix/refactor/chore/prototype}
- **Commit Message**: {message}
- **Scope Score**: {score} (triggered: {list of criteria met})
- **Key Change Areas**: {directories/modules with most impact}

## Scope Boundary

Line-level code quality (naming, formatting, individual function logic, test coverage) is out of scope. Do NOT review these.

## Diff

{same diff size rules as code-reviewer}
```

---

## Error Handling

- **Empty diff** → Report "No changes to review" and stop
- **Git command failure** → Report the error with the failed command
- **Agent failure** → Include partial results from successful agent + note which agent failed
- **Agent reports no issues** → Include "No issues found" in the relevant section

---

## Output Format

**Generate consolidated review summary:**

```markdown
# Code Review Report

**Generated**: {timestamp}
**Scope**: {uncommitted changes / latest commit: {hash}}
**Files Changed**: {count}
**Reviewers Invoked**: code-reviewer{, architect-reviewer if applicable}

---

## 📊 Change Summary

- **Lines**: +{added} -{removed}
- **Files**: {count} ({new_count} new, {modified_count} modified)
- **Directories**: {affected directories}

---

## 🔍 Code Reviewer Feedback

{consolidated output from code-reviewer agent}

---

## 🏗️ Architecture Reviewer Feedback

{if invoked, consolidated output from architect-reviewer agent}
{if not invoked, state: "Architectural review not required for this change scope"}

---

## ✅ Review Checklist

Based on agent feedback, generate action items:

### Critical (Must Fix)

- [ ] {issue 1}
- [ ] {issue 2}

### Warning (Should Fix)

- [ ] {issue 3}

### Suggestion (Consider)

- [ ] {improvement 1}

---

## 📝 Next Steps

{recommended actions based on review results}
```

---

## Important Notes

- **Never modify code directly** - this command only performs review
- **Agent autonomy**: code-reviewer and architect-reviewer may read files, run tests, or analyze dependencies as needed
- **Coding guidelines**: code-reviewer checks for coding standard violations
- **Incremental reviews**: For large changes (>20 files), agents may focus on high-impact areas first
- **Git safety**: All git commands are read-only (status, diff, log, show)

---

## Example Usage

```bash
# Review uncommitted changes
/code-review

# Review specific commit (pass commit hash in conversation)
# User: "Review commit abc123"
# Then: /code-review
```
