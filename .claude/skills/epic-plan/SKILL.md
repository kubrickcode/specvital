---
name: epic-plan
description: Create an execution roadmap from epic analysis. Use after epic-analyze to finalize execution order, define milestones, and prepare workflow-analyze inputs for each sub-task.
disable-model-invocation: true
---

# Epic Planning Command

## User Input

```text
$ARGUMENTS
```

Expected format: `/epic-plan PROJECT-NAME [additional requirements]`

Example:

- `/epic-plan LEGACY-MIGRATION skip validation for task 2`
- `/epic-plan AUTH-SYSTEM prioritize task 3`

You **MUST** consider the user input before proceeding (if not empty).

---

## Outline

1. **Parse User Input**:
   - Extract project name from first word of $ARGUMENTS
   - Extract additional requirements from remaining text
   - If project name missing, ERROR: "Please provide project name: /epic-plan PROJECT-NAME"

2. **Check Prerequisites**:
   - Verify `docs/epic/{project-name}/analysis.md` exists
   - If not exists, ERROR: "Run /epic-analyze first"

3. **Load Analysis**:
   - Extract sub-task list, dependencies, and risks from analysis.md
   - Note user's feedback or adjustments from conversation

4. **Finalize Execution Order**:
   - Confirm or adjust phase structure from analysis.md
   - Incorporate user feedback and additional requirements
   - Define concrete milestones with check conditions

5. **Prepare Task Execution Guides**:
   - For each sub-task, write sufficient context for `/workflow-analyze` in a separate session
   - Include: what to read, what the prior task produced, key constraints
   - This is the **core value** — bridging context between independent sessions

6. **Write Document**:
   - Create `docs/epic/{project-name}/plan.md` (Korean)

---

## Key Rules

### Documentation Language

**CRITICAL**: All documents must be written in **Korean**.

### Core Purpose

**This document bridges context between independent sessions.**
Each sub-task runs in a separate Claude Code session. The plan must provide enough context for each session to start without prior knowledge.

### Must Do

- Reference analysis.md only (no repetition of scope/goals)
- **Execution guide per task**: Sufficient for a cold-start session
- Concrete milestones with check conditions
- Specify what each task should read before starting
- Specify what each task produces for downstream tasks
- Progress tracking checklist

### Must Not Do

- Redefine sub-tasks (they're in analysis.md)
- Repeat dependency graph verbatim (reference it)
- Over-specify implementation (that's workflow-analyze's job)
- Add new sub-tasks not in analysis.md (go back to epic-analyze if needed)
- Estimate time in hours/days

### Execution Guide Principles (CRITICAL)

Each task's execution guide must answer:

1. **What to read first**: Which files/docs to load for context
2. **What prior tasks produced**: Outputs from dependencies that this task needs
3. **Key constraints**: Things that must be respected (from analysis.md or user feedback)
4. **How to start**: The exact `/workflow-analyze` invocation with context
5. **What this task produces**: Outputs that downstream tasks will need

### Parallel Execution (CRITICAL)

**Default is sequential.** Only mark tasks as parallel when ALL conditions are met:

1. **No shared files**: Tasks modify completely different files/modules
2. **No git conflict risk**: Merging both tasks' branches won't produce conflicts
3. **No logical dependency**: Neither task's design decisions affect the other
4. **Truly independent domains**: Different parts of the codebase with no overlap

If in doubt → sequential. The cost of unnecessary sequencing is small (slightly slower). The cost of bad parallelism is high (merge conflicts, rework, inconsistency).

### Milestone Design

- Define milestones at natural **check points**, not at every task boundary
- A milestone = "after these tasks complete, verify this condition"
- Keep milestones few — typically 2-4 for an entire epic
- Each milestone should answer: "Are we still on track, or do we need to re-plan?"

---

## Document Template

File to create: `docs/epic/{project-name}/plan.md` (Korean)

```markdown
# [Project Name] - Execution Plan

> **Analysis**: See `analysis.md`

---

## Execution Order

**Dependencies**: See `analysis.md` > Dependencies section

### Phase 1: [Name]

**Tasks**: [Task names from analysis.md]
**Parallel**: [Yes — all tasks in this phase are independent / No — sequential within phase]
**Goal**: [What's true when this phase completes]

### Phase 2: [Name]

**Tasks**: [Task names]
**Parallel**: [Yes / No / Partial — Task A and B parallel, then Task C]
**Goal**: [What's true when this phase completes]

### Phase 3: [Name]

**Tasks**: [Task names]
**Parallel**: [Yes / No / Partial]
**Goal**: [What's true when this phase completes]

---

## Task Execution Guides

### Task 1: [name]

**Prerequisites**: None / [Task N must be complete]

**Context to load**:

- Read: [specific files or docs]
- Prior output: [what dependency tasks produced — or "N/A" if first task]

**Key constraints**:

- [Constraint from analysis or user feedback]

**Start command**:
```

/workflow-analyze {task-name} [context summary for the session]

```

**This task produces**:

- [What downstream tasks will need from this task's output]

### Task 2: [name]

[Same structure]

### Task N: [name]

[Same structure]

---

## Milestones

### Milestone 1: [Name] (after Phase 1)

**Check condition**: [What to verify]
**If not met**: [What to do — re-plan, adjust scope, etc.]

### Milestone 2: [Name] (after Phase 2)

**Check condition**: [What to verify]
**If not met**: [What to do]

---

## Progress

| Task       | Status | Notes |
| ---------- | ------ | ----- |
| [Task 1]   | -      |       |
| [Task 2]   | -      |       |

Status: `-` pending / `>` in progress / `v` complete / `x` blocked

---

## Quick Reference: Execution Order

> This section is the at-a-glance summary. See details above.

1. [Task name] → `/workflow-analyze {name}`
2. [Task name] → `/workflow-analyze {name}`
3. [Task name] + [Task name] (parallel) → `/workflow-analyze {name}`
4. [Task name] → `/workflow-analyze {name}`

<!-- If all tasks are sequential, omit "(parallel)" markers -->
<!-- If some tasks are parallel, mark only those explicitly -->
```

---

## Execution

Now start the task according to the guidelines above.
