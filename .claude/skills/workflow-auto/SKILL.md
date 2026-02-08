---
name: workflow-auto
description: Execute workflow-analyze → workflow-plan → workflow-execute → commit in sequence without intermediate approval (single commit). Use for small tasks completable in one commit when you want to skip manual phase transitions. NOT for complex features requiring multiple commits, tasks needing validation (/workflow-validate), or when you want to review analysis/plan before execution.
disable-model-invocation: true
---

# Auto Workflow

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

---

## Constraints (vs individual commands)

- No intermediate approval between phases
- Plan limited to **single commit** (no commit splitting, vertical slicing only)
- Automatically generates commit message at the end
- All other processes follow the original commands **exactly** -- no simplification

---

## Document Structure

```
docs/work/{task-name}/
├── analysis.md             (Korean - from workflow-analyze)
├── plan.md                 (Korean - from workflow-plan)
└── summary-commit-1.md     (Korean - from workflow-execute)

commit_message.md           (from commit)
```

---

## Execution

Read and follow each skill file in sequence. You **MUST** read each file to understand the exact process and templates.

1. Read **`.claude/skills/workflow-analyze/SKILL.md`** → execute full process
2. Without stopping → Read **`.claude/skills/workflow-plan/SKILL.md`** → execute (single commit only)
3. Without stopping → Read **`.claude/skills/workflow-execute/SKILL.md`** → execute commit 1
4. Without stopping → Read **`.claude/skills/commit/SKILL.md`** → generate commit_message.md

**Do NOT simplify, skip steps, or use shortened templates.** Follow each skill's exact process.
