---
name: workflow-plan
description: Generate implementation plan with commit-level tasks. Use after workflow-analyze to create detailed, vertically-sliced commit plan.
disable-model-invocation: true
---

# Issue Planning Command

## User Input

```text
$ARGUMENTS
```

Expected format: `/workflow-plan TASK-NAME [additional requirements]`

Example:

- `/workflow-plan REFACTORING limit to 3 commits`
- `/workflow-plan API-REDESIGN keep backward compatibility`

You **MUST** consider the user input before proceeding (if not empty).

---

## Outline

1. **Parse User Input**:
   - Extract task name from first word of $ARGUMENTS (e.g., "REFACTORING")
   - Extract additional requirements from remaining text (e.g., "limit to 3 commits")
   - If task name missing, ERROR: "Please provide task name: /workflow-plan TASK-NAME"

2. **Check Prerequisites**:
   - Verify `docs/work/{task-name}/analysis.md` exists
   - If not exists, use $ARGUMENTS as requirements source and proceed

3. **Load Requirements**:
   - If analysis.md exists: Extract selected approach and completion criteria
   - If not: Interpret user input as requirements (description, URL, file path, etc.)

4. **Identify Impact Scope**:
   - List approximate classes/modules (not specific file names)

5. **Decompose Commits** (Vertical Slicing):
   - Each commit should be independently deployable
   - **Consider additional requirements** from user input (e.g., commit count limits, specific constraints)
   - **Forbid Horizontal Slicing**: Don't separate types/logic/tests/UI into separate commits
   - **Vertical Slicing**: Each commit includes types+logic+tests to provide complete functionality
   - Order: Setup → Core → Integration → Polish
   - Specify verification method and "independently deployable" status for each commit
   - **Test Classification**: Each commit must specify:
     - `TDD`: Logic with conditionals/calculations → Write test FIRST
     - `TEST-AFTER`: UI/Integration/E2E → Implement then test
     - `NO-TEST`: Config/docs/type-only → Skip tests

6. **Review Principle Violations**:
   - Create Complexity Tracking table if coding principle violations are necessary

7. **Write Document**:
   - Create `docs/work/{name}/plan.md` (Korean)

---

## Key Rules

### 📝 Documentation Language

**CRITICAL**: All documents must be written in **Korean**.

### ✅ Must Do

- Checklist-focused
- Reference analysis.md only (no repetition)
- **Vertical Slicing**: Each commit independently deployable
- Reflect coding principles
- Impact scope approximate only

### ❌ Must Not Do

- Redefine problem (it's in analysis.md)
- List specific file names
- Verbose explanations

### 🎯 Vertical Slicing Principles (CRITICAL)

**Each commit must satisfy**:

1. **Build Success**: No compilation errors
2. **Preserve Existing Features**: Pass existing tests
3. **Independently Testable**: Can be tested with this commit alone
4. **Meaningful Value**: Provides real value to users/developers

**❌ Horizontal Slicing Forbidden**:

- Separating types only → logic only → tests only → UI only (X)
- This separation makes each commit functionally incomplete

**✅ Vertical Slicing Example**:

- Commit 1: types + logic + tests + schema (usable with manual config)
- Commit 2: UI integration (complete UX)

### 🧪 Test Classification Rules

**TDD Required** (write failing test first):

- Business logic with conditionals
- Calculations, validations, parsers
- State mutations, data transformations

**TEST-AFTER** (implement then test):

- UI components with visual output
- External API/DB integrations

**NO-TEST** (skip):

- Config files (JSON, YAML, env)
- Documentation, pure types, renames

### 📐 Test Level Selection

**Levels** (lowest effective level first):

| Level  | When to Use                                              |
| ------ | -------------------------------------------------------- |
| `UNIT` | Pure functions, business logic, calculations, validation |
| `INT`  | DB queries, API endpoints, service-to-service            |
| `COMP` | UI components (render, interaction, accessibility)       |
| `E2E`  | Critical user journeys only (auth, payment, core flow)   |

**Decision Quick Reference**:

- Conditional logic? → UNIT
- Database involved? → INT
- UI component? → COMP
- Critical user journey? → E2E (only if below criteria met)

**Mixed Levels**: Use `UNIT + E2E` when feature needs both (e.g., discount logic + checkout flow)

**E2E Justification Required**: E2E tests are easy to write but expensive to run/maintain. Only use E2E when ALL apply:

1. **Cannot test at lower level**: Flow crosses multiple services/pages that mocking defeats the purpose
2. **Revenue/trust critical**: Auth, payment, core conversion paths
3. **Historical bugs**: This flow has broken in production before
4. **Failure is catastrophic**: Regulatory, security, or business-critical

**Anti-pattern**: Writing E2E for scenarios adequately covered by UNIT + INT tests.

### 📊 Phase Structure

- **Phase 1**: Setup
- **Phase 2**: Foundational
- **Phase 3+**: User Stories (we call them Core features)
- **Final Phase**: Polish

→ We apply as Commit order

---

## Document Template

File to create: `docs/work/{task-name}/plan.md` (Korean)

```markdown
# [Task Name] - Implementation Plan

> **Analysis Result**: See `analysis.md`
> **Selected Approach**: [Approach N]

## 📁 Impact Scope (Approximate)

**Main Areas**: [StatusBarManager, ConfigManager, etc.]

---

## 📝 Commit Plan

### ✅ Commit 1: [Title]

**Goal**: [1 sentence - describe complete value provided by this commit]

**Test Requirement**: [TDD | TEST-AFTER | NO-TEST]
**Test Level**: [UNIT | INT | COMP | E2E | UNIT + E2E]

**Task Checklist**:

- [ ] [Specific task 1]
- [ ] [Specific task 2]
- [ ] Test ([LEVEL]): [Specific test case] (if TDD/TEST-AFTER)
- [ ] Build verification: [Build command]

**Verification Method**:

1. [Specific verification method]
2. [User scenario]

**Independently Deployable**: ✅ / ⚠️ [Reason]

---

### ⬜ Commit 2: [Title]

**Goal**: [1 sentence - describe complete value provided by this commit]

**Test Requirement**: [TDD | TEST-AFTER | NO-TEST]
**Test Level**: [UNIT | INT | COMP | E2E | UNIT + E2E]

**Task Checklist**:

- [ ] [Specific task 1]
- [ ] [Specific task 2]
- [ ] Test ([LEVEL]): [Specific test case] (if TDD/TEST-AFTER)
- [ ] Build verification: [Build command]

**Verification Method**:

1. [Specific verification method]
2. [User scenario]

**Independently Deployable**: ✅ / ⚠️ [Reason]

---

## ⚠️ Principle Violation Justification (Only if needed)

| Violation | Why Necessary                | Why Simple Alternative Rejected |
| --------- | ---------------------------- | ------------------------------- |
| Using any | VS Code API types incomplete | Cost to fix @types > benefit    |

---

## 📊 Progress

- [ ] Commit 1
- [ ] Commit 2
```

---

## Execution

Now start the task according to the guidelines above.
