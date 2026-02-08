---
name: handover
description: Generate a comprehensive markdown summary of our conversation for seamless handoff to another AI agent. Use when ending a session or transferring work.
---

# Conversation Handoff Summary

Generate comprehensive summary for AI agent handoff: $ARGUMENTS

## Required Output

This command generates and saves a markdown file. Always:

1. Analyze the conversation
2. Create the file using Write tool
3. Save to: `./handoff-summary-YYYYMMDD-HHMMSS.md`

## Output Template

```markdown
# Project Handoff Summary

Generated: [timestamp]

---

## Overview

Brief description of the project and main objectives

## User Context

- Technical background and preferences
- Specific requirements and constraints
- Communication style preferences

## Original Requirements

- Initial problem statement
- Key goals and success criteria

## Completed Work

### Task 1: [Description]

- What was done and why this approach
- Key code changes and files affected

### Task 2: [Description]

...

## Technical Decisions

### Decision 1: [Topic]

- Options considered, chosen approach, rationale

## Current State

- What's working
- Known issues or limitations

## Pending Tasks

- [ ] Task 1
- [ ] Task 2

## Important Warnings

- Critical information to avoid breaking changes
- Security or performance considerations

## Next Steps

Recommended actions for continuing the work

## Related Files

- File 1: Purpose and recent changes
- File 2: Purpose and recent changes
```

## Writing Guidelines

Handoff summaries are high hallucination risk. Apply these constraints:

- **Distinguish completion status**: "completed" (verified working), "attempted" (may have issues), "discussed" (not implemented)
- **Quote user's exact words** for preferences, constraints, and requirements
- **Only include facts confirmed in conversation** -- do not infer unstated decisions or requirements
- **Mark uncertainty**: If unsure about a detail, write "[needs confirmation]"
- **Capture the user's actual preferences** (e.g., "no emojis", "Korean responses"), not inferred ones

## Self-Verification (Required)

Before saving the file, verify:

- [ ] Every "completed" item has evidence in the conversation
- [ ] No tasks listed as done that were only discussed or attempted
- [ ] User preferences accurately reflect stated preferences, not inferred ones
- [ ] Technical decisions include rationale from actual conversation, not fabricated reasoning
- [ ] Pending tasks distinguish user-requested from agent-suggested
- [ ] No information fabricated to "fill" template sections -- leave empty if nothing to report
