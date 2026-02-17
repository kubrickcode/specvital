# Agent Patterns Reference

Detailed patterns and advanced configurations for agent creation.

## Frontmatter Advanced Fields

### Permission Modes

| Mode                | Behavior                     | Use Case                    |
| ------------------- | ---------------------------- | --------------------------- |
| `default`           | Standard permission prompts  | Normal interactive use      |
| `acceptEdits`       | Auto-accept file edits       | Trusted code generation     |
| `dontAsk`           | Auto-deny permission prompts | Read-only automation        |
| `plan`              | Read-only exploration mode   | Research and planning       |
| `bypassPermissions` | Skip ALL permission checks   | Sandboxed environments only |

### Skills Preloading

Inject skill content into agent context at startup:

```yaml
skills:
  - skill-name-1
  - skill-name-2
```

Use when the agent consistently needs specific domain knowledge that exists as a skill.

### Memory Configuration

Enable persistent memory across sessions:

```yaml
memory: project
```

- `user`: `~/.claude/agent-memory/<name>/` — Personal, all projects
- `project`: `.claude/agent-memory/<name>/` — Project-scoped, shared via VCS
- `local`: `.claude/agent-memory-local/<name>/` — Project-scoped, gitignored

Agent automatically gets Read/Write/Edit for its memory directory. First 200 lines of `MEMORY.md` are injected into system prompt.

### Hooks

Lifecycle hooks scoped to the agent:

```yaml
hooks:
  PreToolUse:
    - matcher: Bash
      command: "scripts/validate_command.sh"
  Stop:
    - command: "scripts/cleanup.sh"
```

Events: `PreToolUse`, `PostToolUse`, `Stop`

## Body Structure Patterns

### Pattern A: Sequential Process

Best for procedural domains (debugging, deployment, migration).

```markdown
You are an expert [ROLE].

When invoked:

1. [Analysis step]
2. [Diagnosis step]
3. [Implementation step]
4. [Verification step]

[Domain-specific guidance]

For each task, provide:

- [Deliverable 1]
- [Deliverable 2]

Focus on [key principle].
```

### Pattern B: Responsibility Framework

Best for architectural/design domains (database, backend, frontend).

```markdown
You are a [ROLE] with expertise in [DOMAIN].

## Core Responsibilities

### [Area 1]

- [Specific guidance]

### [Area 2]

- [Patterns and approaches]

## Architecture Patterns

[Decision tables, framework comparisons]

## Implementation Workflow

1. [Phase 1]
2. [Phase 2]

## Tool Selection

- **[Tool]**: [Purpose and rationale]

## Key Principles

- [Principle 1]
- [Principle 2]
```

### Pattern C: Decision Framework

Best for evaluative/analytical domains (tech advisor, product strategist).

```markdown
You are a [ROLE] specializing in [DOMAIN].

## Analysis Process

1. [Scope detection]
2. [Data gathering]
3. [Framework application]
4. [Synthesis]

## Decision Matrix

| Criterion | Weight | [Option A] | [Option B] |
| --------- | ------ | ---------- | ---------- |

## Output Format

### [Report structure]

## Key Principles

- [Evaluation standards]
```

### Pattern D: Research + Synthesis

Best for information-gathering domains (research, competitive analysis).

```markdown
You are a [ROLE] expert at [CAPABILITY].

## Focus Areas

- [Area 1]
- [Area 2]

## Research Process

1. [Query design]
2. [Search execution]
3. [Verification]
4. [Synthesis]

## Output Format

### Key Findings

- [Finding with source]

### Source Assessment

| Source | Credibility | Notes |

### Recommendations

- [Actionable insight]
```

### Pattern E: Tool-Focused Optimization

Best for meta/system optimization domains.

```markdown
You are a [ROLE] specializing in [SYSTEM].

## Core Workflow

1. [Analyze]
2. [Identify issues]
3. [Present options]
4. [Implement]
5. [Verify]

## [Domain] Expertise

[Specific knowledge areas]

## Diagnostic Process

[Issue categorization and measurement]

## Optimization Strategies

[Concrete improvement techniques]

## Critical Requirements

- [Must-do items]

## Common Pitfalls

- [What to avoid]
```

## Quality Checklist

Validate every agent against:

- [ ] `description` is under 200 characters with "Use PROACTIVELY" trigger
- [ ] `tools` explicitly listed (not relying on inheritance)
- [ ] Opening line establishes expertise: `You are a [ROLE]...`
- [ ] Contains numbered workflow steps
- [ ] Has "Output Format" section (if agent produces structured output)
- [ ] Has "Key Principles" or equivalent closing section
- [ ] Uses H2/H3 hierarchy only (no H4)
- [ ] No redundant content that Claude already knows
- [ ] Closing statement is a principle or success criteria
- [ ] File size matches domain complexity (30-300 lines)

## Tool Selection Reference

Common tool combinations by agent type:

| Agent Type        | Recommended Tools                                              |
| ----------------- | -------------------------------------------------------------- |
| Code reviewer     | Read, Grep, Bash                                               |
| Implementation    | Read, Write, Edit, Bash, Glob, Grep                            |
| Research          | Read, WebFetch, WebSearch                                      |
| Architecture      | Read, Write, Edit, Bash, Glob, Grep                            |
| Documentation     | Read, Write, Edit, Glob, Grep                                  |
| Meta/optimization | Read, Write, Edit, Bash, Glob, Grep, WebFetch, AskUserQuestion |
| Context/planning  | Read, Write, Edit, TodoWrite                                   |

**Collaboration delegation pattern:**

```markdown
## Tool Selection

Essential tools:

- **Read/Grep/Glob**: Codebase analysis
- **Edit/Write**: Implementation
- **Bash**: Build verification

Collaboration:

- **prompt-engineer**: Complex prompt optimization
- **security-auditor**: Security review
- **test-expert**: Test strategy design
```

## Anti-Patterns to Avoid

- **Scope creep**: Agent tries to do everything in a domain
- **Tool over-granting**: Omitting `tools` field (inherits all)
- **Vague triggers**: Description says "helps with X" without specifics
- **Redundant content**: Explaining things Claude already knows
- **Missing output format**: Agent produces inconsistent results
- **Wrong model**: Using opus for simple search tasks
- **No workflow**: Body is a wall of text without numbered steps
