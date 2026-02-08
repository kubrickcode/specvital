---
name: ultra-think
description: Deep analysis and problem solving with multi-dimensional thinking. Use when facing complex architectural decisions, strategic planning, or problems requiring thorough analysis from multiple perspectives.
---

# Deep Analysis Mode

Analyze the following problem with maximum depth and rigor: $ARGUMENTS

## Analysis Requirements

1. **Identify the core challenge** — separate symptoms from root causes
2. **Surface hidden constraints** — question stated assumptions, identify unstated ones
3. **Generate 3+ distinct approaches** — genuinely different strategies, not variations of one idea
4. **Evaluate trade-offs honestly** — not all pros/cons are equal; weight by actual impact
5. **Recommend with rationale** — state confidence level and what would change your mind

## Quality Standards

### Anti-Overconfidence Rules

- **Default confidence is Medium.** High confidence requires: one approach is Pareto-dominant, OR a hard constraint eliminates alternatives.
- **State uncertainty explicitly**: "I'm uncertain about X because Y" is more valuable than false precision.
- **No manufactured balance**: If one option is clearly superior, say so. If options are genuinely close, say that instead.
- **Forbidden language**: "obviously", "clearly", "simply", "no-brainer", "the only real option"
- **Distinguish fact from inference**: Mark claims as "confirmed" (from code/docs), "likely" (strong evidence), or "uncertain" (inference).

### Depth Over Breadth

- Fewer perspectives analyzed deeply beat many perspectives analyzed superficially.
- Skip perspectives irrelevant to the specific problem (not every problem needs a "business perspective").
- Second-order effects matter more than comprehensive first-order coverage.

## Output Structure

```markdown
## Problem Analysis

- Core challenge (1-2 sentences)
- Key constraints
- What makes this decision difficult

## Approaches

### Approach 1: [Name]

- Description, pros, cons, risks
- Confidence: [High/Medium/Low] — [basis]

### Approach 2: [Name]

[Same structure]

## Recommendation

- **Recommended**: [Approach N]
- **Confidence**: [High/Medium/Low]
- **Rationale**: [Why — acknowledging what you're giving up]
- **Would change if**: [What new information would alter this recommendation]

## Uncertainties

- [What you don't know that matters]
- [Where you might be wrong]
```

## Self-Verification (Required)

Before presenting the analysis, verify:

- [ ] Each approach has at least one genuine weakness identified (not a token con)
- [ ] Confidence level is stated and justified for each approach and the final recommendation
- [ ] "Uncertainties" section is non-empty and substantive
- [ ] No recommendation uses forbidden confident language
- [ ] Approaches are genuinely distinct (not minor variations)
- [ ] Analysis depth matches problem complexity (simple problems get concise analysis)
