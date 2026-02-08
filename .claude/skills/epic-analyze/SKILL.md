---
name: epic-analyze
description: Decompose a large-scale project into manageable sub-tasks with dependency analysis. Use when a task is too large for a single workflow cycle — requires splitting into multiple independent work units before implementation.
disable-model-invocation: true
---

# Epic Analysis Command

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

---

## Outline

1. **Parse User Input**:
   - Extract project description from $ARGUMENTS
   - Generate project name (2-4 words, English, hyphen-separated)

2. **Understand the Whole**:
   - Define project goal in 1-2 sentences
   - Identify the boundary: what's in, what's explicitly out
   - Understand current state (existing code, infrastructure, constraints)

3. **Identify Sub-Tasks**:
   - Break the project into the **minimum number** of independent work units
   - Each sub-task should be completable within a single workflow cycle (analyze → plan → execute)
   - Name each sub-task (2-4 words, English, hyphen-separated)

4. **Analyze Dependencies**:
   - Map which sub-tasks block others
   - Identify what can run in parallel
   - Find the critical path

5. **Assess Risks**:
   - Technical uncertainties requiring validation (`/workflow-validate`)
   - External dependencies or unknowns
   - Areas where the decomposition itself is uncertain

6. **ADR Assessment**:
   - After decomposition, evaluate whether any project-level decisions meet ADR criteria (PRICE)
   - Epic-level decisions often qualify: technology choices, architectural patterns, cross-cutting concerns
   - If any PRICE criterion applies, add the ADR assessment section to the document

7. **Assess Decomposition Confidence**:
   - Apply Decision Honesty rule (same as workflow-analyze)
   - **Confident**: Clear boundaries, well-understood domain
   - **Tentative**: Some boundaries may shift after early tasks reveal more

8. **Write Document**:
   - Create `docs/epic/{project-name}/analysis.md` (Korean)

---

## Key Rules

### Documentation Language

**CRITICAL**: All documents must be written in **Korean**.

### Thinking Order

**Understand the whole project first. Decompose second.
Premature decomposition without understanding the whole leads to wrong boundaries.**

### Decomposition Principles (CRITICAL)

1. **Minimum viable split**: Fewest sub-tasks that cover the full scope. If two tasks always change together, they're one task.
2. **Session-independent**: Each sub-task must be executable in a separate Claude Code session with only its own context (analysis.md, plan.md) — no implicit dependency on "what the previous session knew."
3. **Vertical, not horizontal**: Don't split by layer (DB → API → UI). Split by capability or domain boundary.
4. **Size calibration**:
   - Each sub-task = 1-5 commits when planned with `/workflow-plan`
   - If a sub-task feels like 1 commit → it's too small, merge with neighbors
   - If a sub-task feels like 10+ commits → it needs further splitting
5. **Workflow-ready output**: Each sub-task description must be sufficient as input to `/workflow-analyze {task-name}`

### Anti-Bias Rules (CRITICAL)

**Core: Don't over-engineer the decomposition itself. The goal is clarity, not perfection.**

1. **Over-decomposition bias**
   - AI tends to split too much ("just to be safe")
   - Ask: "Would merging these two tasks lose anything meaningful?"
   - If the answer is no → merge them

2. **Scope creep bias**
   - AI tends to add "nice to have" sub-tasks not in the original request
   - Every sub-task must trace back to the user's stated goal
   - If it doesn't → exclude it or mark it explicitly as "future consideration"

3. **Perfectionism bias**
   - AI tends to add "setup" and "cleanup" tasks that aren't needed
   - Infrastructure/tooling tasks are only valid if they unblock other tasks
   - Don't create sub-tasks for things the team already knows how to do

4. **Symmetry bias**
   - AI tends to make all sub-tasks similar in size
   - Uneven size distribution is natural — don't force uniformity

5. **Decision Honesty** (same principle as workflow-analyze)
   - **Confident**: Domain is well-understood, boundaries are clear → state definitively
   - **Tentative**: Some boundaries may shift → mark uncertain areas, recommend revisiting after Phase 1
   - **Default is Tentative.** Confident requires well-understood domain + clear boundaries.

### Must Do

- Project goal in 1-2 sentences
- **Minimum** number of sub-tasks (prefer fewer, larger tasks)
- Clear dependency graph
- Each sub-task named and described enough for `/workflow-analyze` input
- Estimate relative size (S/M/L) per sub-task
- Identify what needs `/workflow-validate` before proceeding
- Mark uncertain boundaries explicitly

### Must Not Do

- **Create sub-tasks for obvious setup** (project init, linter config, etc.) unless genuinely complex
- **Split by technical layer** (types → logic → tests → UI)
- **Add tasks not in the original scope** without explicit marking
- Estimate time in hours/days (use relative size only)
- Over-specify implementation details (that's workflow-analyze's job)

---

## Document Template

File to create: `docs/epic/{project-name}/analysis.md` (Korean)

```markdown
# [Project Name] - Epic Analysis

## Overview

**Project Goal**: [1-2 sentences — what this project achieves when complete]
**Trigger**: [Why now? What prompted this project?]

---

## Scope

**Included**:

- [What's in scope]

**Excluded**:

- [What's explicitly out — and why]

**Current State**: [Brief description of what exists today]

---

## Sub-Task Decomposition

<!-- Confident: use this format -->

**Confidence**: Confident (boundaries are clear)

<!-- Tentative: use this format instead -->

**Confidence**: Tentative (boundaries may shift)
**Uncertain Areas**: [Which sub-task boundaries might change and why]
**Revisit After**: [Which task(s) should complete before re-evaluating]

### Task 1: [name] (Size: S/M/L)

**Goal**: [1 sentence — what this task achieves independently]
**Why separate**: [Why this isn't part of another task]
**Workflow input**: [Key context for `/workflow-analyze {name}` — 2-3 sentences max]

### Task 2: [name] (Size: S/M/L)

[Same structure]

### Task N: [name] (Size: S/M/L)

[Same structure]

---

## Dependencies
```

Task 1 ──→ Task 3 ──→ Task 5
Task 2 ──→ Task 4 ──┘

```

**Critical Path**: Task 1 → Task 3 → Task 5
**Parallelizable**: Task 1 and Task 2 can proceed simultaneously

---

## Risks & Validation Needs

**Needs `/workflow-validate`**:

- [Technical uncertainty 1]: [What to validate and why]

**External Dependencies**:

- [Dependency]: [Impact if unavailable]

---

## Suggested Execution Order

**Phase 1**: [Task names — foundation]
**Phase 2**: [Task names — core features]
**Phase 3**: [Task names — integration/polish]

> Start with: `/workflow-analyze {task-1-name}`

---

## ADR Assessment

<!-- Include this section ONLY when project-level decisions meet PRICE criteria -->
<!-- Omit entirely if no criterion applies -->

**Decisions requiring ADR**:

- [Decision 1]: PRICE criteria [R, I] — [brief reason]
- [Decision 2]: PRICE criteria [P] — [brief reason]

**Recommendation**: Create ADRs before starting sub-task execution
**Suggested Timing**: After epic analysis approval, before first `/workflow-analyze`

> To create: `/adr [topic]`

---

## Self-Check

- [ ] Minimum number of sub-tasks (couldn't merge any two without losing independence)
- [ ] Every sub-task traces to the stated project goal
- [ ] No "nice to have" tasks smuggled in
- [ ] Each sub-task is session-independent (no implicit context dependency)
- [ ] Vertical split, not horizontal (each task delivers a complete capability)
- [ ] Decomposition confidence honestly assessed
```

---

## Execution

Now start the task according to the guidelines above.
