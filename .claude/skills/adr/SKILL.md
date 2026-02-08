---
name: adr
description: Create or update Architecture Decision Records (ADR) in docs/decisions/. Use when recording architecturally significant decisions — triggered by PRICE criteria (Policy, Reversibility, Impact, Constraint, Exception), from analysis phase recommendations (workflow-analyze, epic-analyze), or manual invocation during any conversation.
disable-model-invocation: true
---

# ADR Creation Command

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

---

## Outline

1. **Parse User Input**:
   - Extract decision topic from $ARGUMENTS
   - Generate ADR slug (2-4 words, English, hyphen-separated)

2. **Gather Context**:
   - Use conversation history and any user-tagged files as context
   - If analysis.md is tagged: extract decision, options, and rationale from it (distill, don't duplicate)

3. **Scan Existing ADRs**:
   - List file names in `docs/decisions/` to determine next number and identify potential duplicates by slug
   - Read only candidate files (similar slug) to check `status`, `name`, `description` before superseding

4. **Write ADR Document**:
   - Create `docs/decisions/NNNN-slug.md` (Korean)
   - If derived from analysis.md: distill (don't duplicate) — extract only the decision and rationale
   - If standalone: write from conversation context

5. **Cross-Link** (if applicable):
   - If superseding an existing ADR: update old ADR's status

6. **Report**:
   - Confirm file created
   - Show PRICE criteria that triggered the ADR

---

## Key Rules

### Documentation Language

**CRITICAL**: ADR documents must be written in **Korean**.
Technical terms may include English in parentheses.

### Distill, Don't Duplicate

When source material (analysis document, conversation) is available:

- Source = detailed analysis or discussion (task-scoped)
- ADR = distilled decision record (0.5-1 page, project-scoped)
- Extract: what was decided, why, what was given up
- Do NOT copy the full source into the ADR

### Numbering

- Zero-padded 4 digits: `0001`, `0002`, ...
- Scan `docs/decisions/` for the highest existing number
- If directory doesn't exist, create it and start at `0001`

### Status Values

- `수락됨` (Accepted) — default for new ADRs
- `대체됨` (Superseded) — use with `supersedes` field
- `폐기됨` (Deprecated)

---

## Must Do

- Include concrete context (not abstract descriptions)
- List all considered options with pros/cons
- Explicitly state accepted trade-offs
- Link to version-controlled sources (PR, commit, other ADRs) when available

## Must Not Do

- Verbose explanations — keep under 1 page
- Duplicate source material wholesale
- Omit trade-offs or pretend the decision has no downsides
- Create ADR for trivial decisions (see `rules/adr.md` exclusion list)

---

## Document Template

File to create: `docs/decisions/NNNN-slug.md` (Korean)

```markdown
---
name: [ADR 제목]
description: [결정 사항 1줄 요약]
status: 수락됨
date: YYYY-MM-DD
supersedes: null
---

# ADR-NNNN: [제목]

## 맥락

[이 결정이 필요한 상황. 구체적 시나리오 포함.]

## 결정 요인

- [핵심 요인 1]
- [핵심 요인 2]

## 검토한 선택지

### 선택지 1: [이름]

- 장점: ...
- 단점: ...

### 선택지 2: [이름]

- 장점: ...
- 단점: ...

## 결정

[선택한 접근법]. [핵심 근거 1-2문장.]

## 결과

- 긍정적: ...
- 수용한 트레이드오프: ...

## 관련 문서

- [버전 관리되는 출처 — PR, 커밋, 다른 ADR 등]
```

---

## Execution

Now start the task according to the guidelines above.
