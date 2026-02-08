---
name: blueprint
description: Transform verbose natural language requests into structured bilingual documentation (Korean for review + English for AI prompts). Use when you need to clarify and structure a complex request before implementation.
---

# Generate Request Blueprint

Transform unstructured natural language request into clear, systematic bilingual documentation: $ARGUMENTS

This command works for ALL types of requests -- technical, non-technical, decision-making, planning, personal, creative, etc.

## Core Principles

1. **Simple request -> Simple document**: 3-4 sections sufficient
2. **Complex request -> Detailed document**: Add only what's needed
3. **Non-technical request**: Skip technical environment, dependencies
4. **Core goal**: Reduce verbosity, structure essentials only

## Language Policy

The English version must include this notice at the top:

> **Language Requirement: All responses, conversations, and outputs should be in Korean.**

This ensures AI agents respond in Korean even when receiving English prompts.

## Document Structure

Create markdown file in **root directory** with filename: `blueprint-[brief-title]-[timestamp].md`

```markdown
# [Request Title]

---

## Korean Version (For Review)

> Version for user to review and confirm request content

[Include only necessary sections based on context]

### Possible Sections (use selectively)

- Request Overview
- Purpose & Background
- Specific Requirements
- Constraints & Warnings
- Expected Outcomes
- Technical Environment (technical requests only)
- Priority (if needed)
- References (if available)
- Additional Context (if needed)

---

## English Version (For AI Prompt)

> **Language Requirement: All responses, conversations, and outputs should be in Korean.**
>
> This section is optimized for AI agent consumption.

[Same structure as Korean version, written in English]
[Include same sections only, do not add unnecessary ones]

---

## Next Steps

1. Review the Korean version
2. Modify if needed
3. Share the English version with the next AI agent
```

## Writing Guidelines

Document only what the user explicitly stated. This prevents hallucination and ensures the blueprint accurately represents user intent.

### DO:

- **Record user's exact words** for ambiguous expressions
- **Keep vague terms vague** -- do not interpret
- **Only include explicitly mentioned** requirements, constraints, and context
- **Preserve user's original phrasing** when unclear
- **Mark unclear points** as "Needs clarification: ..." rather than guessing

### DON'T:

- **Extract "implicit" requirements** -- explicit only
- **Interpret vague expressions** (e.g., "tight circumstances" as specific financial problems)
- **Add details** the user didn't mention
- **Expand abbreviations** or casual language into formal assumptions
- **Infer user's situation** beyond what they stated

## Self-Verification (Required)

Before saving the file, verify:

- [ ] Only user's explicit statements included (no inferred requirements)
- [ ] Ambiguous expressions preserved verbatim with clarification markers
- [ ] Korean and English versions contain identical scope (no section drift)
- [ ] Section count matches request complexity (no forced empty sections)
- [ ] No "implicit requirements" fabricated to fill the document
- [ ] Language policy notice present in English version
