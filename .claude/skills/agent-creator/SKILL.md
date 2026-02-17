---
name: agent-creator
description: Guide for creating effective Claude Code custom agents (.claude/agents/ markdown files). Use when users want to create a new agent or update an existing agent that provides specialized domain expertise as a subagent. Covers frontmatter configuration, system prompt design, tool scoping, and model selection.
---

# Agent Creator

Guide for creating effective Claude Code custom agents.

## About Agents

Agents are specialized subagent definitions (`.claude/agents/*.md`) that provide domain expertise. When invoked, the agent's markdown body becomes the subagent's system prompt.

**Key distinction from other config types:**

- **Skills**: Automate repetitive multi-step workflows with bundled resources
- **Rules**: Enforce coding standards and constraints (always loaded)
- **Agents**: Define domain expert personas with scoped tools and models

## Core Principles

### Description is King

The `description` field is the primary auto-triggering mechanism. Claude reads ALL agent descriptions to decide delegation. Write action-oriented triggers:

- Good: `"Database performance optimization specialist. Use PROACTIVELY for slow queries, indexing strategies, and execution plan analysis."`
- Bad: `"Helps with databases"`

### Single Responsibility

One agent, one domain. If an agent needs "and" in its purpose, split it.

### Principle of Least Privilege

Always specify `tools` explicitly. Never rely on tool inheritance (inherits ALL tools if omitted). Grant only what the agent needs.

### Token Economy

The agent body is injected as a system prompt every invocation. Keep it lean — include only core workflow and role definition in the body.

## Agent Creation Process

1. Understand the domain with concrete examples
2. Configure frontmatter fields
3. Initialize the agent (run init_agent.py)
4. Write the system prompt body
5. Validate the agent (run validate_agent.py)
6. Iterate based on real usage

### Step 1: Understand the Domain

Identify:

- What specific tasks should this agent handle?
- What would a user say that should trigger this agent?
- What tools does the agent need?
- How complex is the domain? (determines body size)

Conclude when there is a clear sense of scope and triggers.

### Step 2: Configure Frontmatter

Required fields:

| Field         | Purpose                                                          |
| ------------- | ---------------------------------------------------------------- |
| `name`        | kebab-case identifier matching filename                          |
| `description` | Auto-triggering text. Include "Use PROACTIVELY when..." triggers |

Optional fields:

| Field             | Purpose                                       | Default   |
| ----------------- | --------------------------------------------- | --------- |
| `tools`           | Allowlist of tools                            | All tools |
| `disallowedTools` | Denylist of tools                             | None      |
| `model`           | `haiku`, `sonnet`, `opus`, `inherit`          | `inherit` |
| `maxTurns`        | Max agentic turns                             | Unlimited |
| `skills`          | Skills to preload at startup (see note below) | None      |
| `mcpServers`      | MCP servers available                         | None      |
| `permissionMode`  | Permission behavior                           | `default` |
| `memory`          | `user`, `project`, or `local`                 | None      |

**Model selection guide:**

- `haiku`: Exploration, search, simple tasks (cost-efficient)
- `sonnet`: Code review, implementation, balanced tasks
- `opus`: Architecture decisions, complex reasoning
- `inherit`: Match user's session model (recommended default)

**Tool selection common sets:**

- Read-only research: `Read, Grep, Glob`
- Research + web: `Read, Grep, Glob, WebFetch, WebSearch`
- Implementation: `Read, Write, Edit, Bash, Glob, Grep`
- Full stack: `Read, Write, Edit, Bash, Glob, Grep, WebFetch, WebSearch`

**Note on `skills` field:**

The `skills` field injects skill SKILL.md content into the subagent's system prompt at spawn time. Subagents are isolated — they cannot see the parent conversation's context, so this field provides domain knowledge they would otherwise lack.

This is only relevant when the agent runs as an **auto-delegated subagent** (Claude spawns it via Task tool) or in **`claude --agent` / Agent Teams** mode. For manual invocation workflows (e.g., `@agent-name` in conversation), the main conversation already has full context, making `skills` unnecessary.

For details on all frontmatter fields and advanced patterns, see [references/agent-patterns.md](references/agent-patterns.md).

### Step 3: Initialize the Agent

Run the initialization script to create a template agent:

```bash
scripts/init_agent.py <agent-name> --path <output-directory>
```

The script creates a `.md` file with proper frontmatter and TODO placeholders.

Skip this step if updating an existing agent.

### Step 4: Write the System Prompt Body

The markdown body after frontmatter becomes the subagent's system prompt. Structure it based on domain complexity:

**Minimal (30-40 lines)** — Narrow, focused domains:

```markdown
You are an expert [ROLE] specializing in [DOMAIN].

When invoked:

1. [Step 1]
2. [Step 2]
3. [Step 3]

For each task, provide:

- [Deliverable 1]
- [Deliverable 2]

Focus on [key principle].
```

**Standard (80-140 lines)** — Most agents:

```markdown
You are a [ROLE] with expertise in [DOMAIN].

## Core Workflow

When invoked:

1. [Step 1]
2. [Step 2]
3. [Step 3]

## [Domain Area 1]

- [Specific guidance]
- [Decision criteria]

## [Domain Area 2]

- [Patterns and approaches]

## Tool Selection

Essential tools:

- **Read/Grep**: [Purpose]
- **Edit/Write**: [Purpose]

Collaboration:

- **[other-agent]**: [When to delegate]

## Output Format

[Template or example of expected output]

## Key Principles

- [Principle 1]
- [Principle 2]
```

**Comprehensive (150+ lines)** — Complex domains with multiple phases:

- Add decision tables/matrices
- Add phased workflows
- Add anti-patterns to avoid
- Add success criteria

For detailed body structure templates, see [references/agent-patterns.md](references/agent-patterns.md).

**Writing guidelines:**

- First line: `You are a [ROLE]...` — establishes expertise
- Use imperative form for instructions
- Numbered steps for sequential processes
- Bullet points for options/checklists
- Tables for decision frameworks
- H2 for main sections, H3 for subsections (never H4)
- Include "Tool Selection" section with rationale
- End with a key principle or success criteria

### Step 5: Validate the Agent

Run the validation script:

```bash
scripts/validate_agent.py <path/to/agent.md>
```

Checks:

- YAML frontmatter format and required fields
- Description quality and length
- Tool names validity
- Model value validity
- Body structure conventions

### Step 6: Iterate

1. Test by delegating real tasks to the agent
2. Observe where it struggles or produces unexpected output
3. Refine the system prompt, tool scope, or model selection
4. Re-validate after changes

## Bilingual Projects

When the project has both `.claude/agents/` and `.claude.ko/agents/`:

- Create both versions simultaneously
- English version: `.claude/agents/<name>.md`
- Korean version: `.claude.ko/agents/<name>.md`
- Maintain identical structure and logic, only translate content
