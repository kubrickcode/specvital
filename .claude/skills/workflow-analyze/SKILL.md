---
name: workflow-analyze
description: Analyze issue and generate analysis.md with solution approaches. Use when starting a new task or feature to thoroughly analyze the problem, compare solution options, and produce a structured decision document. Triggers on task analysis, solution comparison, and pre-implementation investigation.
disable-model-invocation: true
---

# Issue Analysis Command

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

---

## Outline

1. **Parse User Input**:
   - Extract task description from $ARGUMENTS
   - Generate task name (2-4 words, English, hyphen-separated)

2. **Define Problem**:
   - Identify current situation, problem to solve, and goals

3. **Investigate Solutions**:
   - Research 2-4 approaches
   - Analyze pros/cons of each approach
   - Evaluate based on coding principles

4. **Assess Decision Confidence & Final Selection**:
   - After completing ALL analysis, assess decision confidence (see Decision Honesty rule)
   - **Decisive** → Select the approach, clarify rejection reasons
   - **Advisory** → Recommend the approach, present alternatives with deciding factors for user

5. **Generate Questions**:
   - Select **maximum 3** unclear items only
   - Use reasonable defaults + Assumptions for the rest

6. **ADR Assessment**:
   - After completing the analysis, evaluate whether the decision meets ADR criteria (PRICE)
   - If any PRICE criterion applies, add the ADR assessment section to the document

7. **Write Document**:
   - Create `docs/work/{task-name}/analysis.md` (Korean)

---

## Key Rules

### 📝 Documentation Language

**CRITICAL**: All documents must be written in **Korean**.

### 📝 Writing Order

**Complete the analysis body (Problem Definition through Solution Investigation) in your reasoning first.
Write the Conclusion (TL;DR) section last, even though it appears first in the document.**

### 🎯 Balanced Analysis Principles

1. **Problem Definition**: Include concrete scenarios rather than abstract descriptions
2. **Solution Approach**: Consider both user perspective and technical feasibility
3. **Implementation Method**: Keep core changes concise, specify new dependencies
4. **Completion Criteria**: Separate feature verification and technical implementation
5. **Weighted Evaluation**: Not all pros/cons are equal. Identify which criteria matter most for this problem before comparing. A single critical advantage can outweigh multiple minor ones.

### 🎭 Anti-Bias Rules (CRITICAL)

**Core Principle: Be honest. Don't hide what exists, don't invent what doesn't.**

1. **Analyze First, Conclude Later**
   - Complete ALL approach analysis BEFORE making any selection

2. **Write What's True**
   - If A has 5 pros and B has 1, write exactly that (no forced balance)
   - If side effects exist, state them; if none, say so (no hiding, no inventing)
   - When rejecting an approach, acknowledge what you're giving up

3. **Neutral Language**
   - Forbidden: "simply", "obviously", "clearly", "just", "merely"
   - Use: "may", "could", "potentially"

4. **Decision Honesty (CRITICAL)**
   - After completing all analysis, assess decision confidence before writing the conclusion
   - **Decisive** (select definitively) — ALL conditions must be met:
     - One approach is superior or equal in all evaluation dimensions, inferior in none (Pareto dominant)
     - OR a technical/project constraint eliminates other approaches
     - OR an explicit, unambiguous project convention determines the answer
     - AND the deciding factor does NOT depend on user/team preferences or business context unknown to you
   - **Advisory** (recommend, let user decide) — ANY condition is sufficient:
     - Approaches have comparable trade-off profiles with strengths in different dimensions (name the dimensions)
     - The deciding factor is a subjective preference (e.g., simplicity vs flexibility)
     - The choice requires business/product context you do not have
     - Risk tolerance or time horizon is the differentiator
     - Differences between approaches are marginal
   - **Default is Advisory.** Decisive requires proving no legitimate trade-off remains.
   - Do NOT manufacture certainty. If the honest answer is "these are close", say so.

   **Decision Type Examples**:
   - Decisive: Project uses PostgreSQL; Approach A is PostgreSQL-based, Approach B is MongoDB-only → technical constraint, Decisive
   - Advisory: Approach A is simpler but less extensible; Approach B is complex but highly extensible → depends on future scale and team preference → Advisory

   **Additional forbidden expressions when using Decisive**:
   - "no need to bother with ~" (dismissive of alternatives)
   - "practically the only choice" (unless literally one option remains)
   - "other approaches are not realistic" (requires specific evidence)

### ✅ Must Do

- Conclusion first (Executive Summary)
- **Include concrete scenarios** (Problem Definition) - 1-4 cases depending on issue
- **Explain user impact** (Solution Approach) - How will it be used?
- **Evaluate technical feasibility** (Solution Approach) - Why is it suitable?
- Evaluate solutions based on coding principles
- **NEEDS CLARIFICATION maximum 3** - use reasonable defaults for the rest

### ❌ Must Not Do

- **Only abstract problem definitions** (without concrete examples)
- **Excessive technical details** (going into implementation level)
- **Missing user perspective** (talking only about technology)
- **Listing obvious things** (e.g., React state, JSON schema)

---

## Document Template

File to create: `docs/work/{task-name}/analysis.md` (Korean)

```markdown
# [Task Name] - Analysis Result

## 🎯 Conclusion (TL;DR)

<!-- Decisive: use this format -->

**Decision Type**: Decisive
**Selected Approach**: [Approach Name]
**Key Reason**: [Why this approach was chosen - 1-2 sentences]

<!-- Advisory: use this format instead -->

**Decision Type**: Advisory (user confirmation required)
**Recommended Approach**: [Approach Name]
**Key Reason**: [Why this approach is recommended - 1-2 sentences]
**Consider Alternative When**: [Under what conditions Approach M would be better]

---

## 📋 Problem Definition

**Concrete Scenarios** (as needed):

- [Case 1: What problem occurs in what situation]
- [Case 2: Another specific situation]

**Current Problem**: [What is inconvenient or impossible]
**Goal**: [What this task aims to achieve]

---

## 📦 Scope

**Included**:

- [Scope Item 1]
- [Scope Item 2]

**Excluded**:

- [What to explicitly exclude]

---

## 🔍 Solution Investigation

### Approach 1: [Name]

**Method**:

- [How to solve - including user impact and technical method]

**Pros**:

- [List ALL genuine pros]

**Cons**:

- [List ALL genuine cons]

**Side Effects/Risks** (if any):

- [What negative consequences may occur]

### Approach 2: [Name]

[Same structure]

---

## ✅ Final Selection

<!-- Decisive: use this format -->

**Decision Type**: Decisive
**Adopted**: Approach N
**Decisive Basis**: [Which specific criterion was met — technical constraint / Pareto dominance / project convention]

**Trade-offs Accepted** (if any):

- [Benefit from other approaches that won't be available]

**Known Risks** (if any):

- [What could go wrong]

**Selection Reason**:

- [Why this approach is most suitable]

**Rejected Approaches** (include what you're losing):

- Approach X: [Rejection reason] — Would have provided: [lost benefit]

<!-- Advisory: use this format instead -->

**Decision Type**: Advisory (user confirmation required)
**Recommended**: Approach N
**Strong Alternative**: Approach M

**Comparison on Deciding Factors**:

| Factor                 | Approach N | Approach M |
| ---------------------- | ---------- | ---------- |
| [Key differentiator 1] | ...        | ...        |
| [Key differentiator 2] | ...        | ...        |

**Why Approach N is Recommended**:

- [Reasoning - acknowledging it's a preference call]

**When to Choose Approach M Instead**:

- [Specific conditions under which the alternative is better]

**Decision Guide**:

- Prioritize [Priority A] → Approach N
- Prioritize [Priority B] → Approach M

---

## 🛠️ Implementation Method

**Core Changes**:

- [Change Item 1]
- [Change Item 2]

**New Dependencies**: [None / Name]

---

## 🎯 Completion Criteria

**Feature Verification**:

- [ ] [Scenario 1 verification]
- [ ] [Existing feature regression test]

**Technical Implementation**:

- [ ] [Core change implementation]
- [ ] [Unit test writing]

---

## ❓ Needs Confirmation

**Current Assumptions**:

- [Assumption 1]: [Default value]

**If Needed**:

- [Items needing additional confirmation - maximum 3]

---

## 📌 ADR Assessment

<!-- Include this section ONLY when at least one PRICE criterion applies -->
<!-- Omit entirely if no criterion applies -->

**PRICE Criteria Met**: [P/R/I/C/E — list which apply with brief reason]
**Recommendation**: Record as ADR before proceeding to planning phase
**Suggested ADR Title**: `NNNN-short-description`

> To create: `/adr [topic]`

---

## 🔍 Objectivity Self-Check

- [ ] Wrote what's true (no hiding, no inventing)
- [ ] Rejected approaches' lost benefits acknowledged
- [ ] Used neutral language ("may/could", not "will/obviously")
- [ ] Decision confidence honestly assessed (did not force certainty on a subjective choice)
- [ ] **If Decisive**: Would a colleague with a different opinion still find this convincing?
- [ ] **If Decisive**: Were alternative approaches' strengths not understated?
```

---

## Execution

Now start the task according to the guidelines above.
