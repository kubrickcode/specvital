#!/usr/bin/env python3
"""
Agent Initializer - Creates a new agent from template

Usage:
    init_agent.py <agent-name> --path <output-directory>

Examples:
    init_agent.py my-specialist --path .claude/agents
    init_agent.py data-analyst --path .claude.ko/agents
"""

import sys
import re
from pathlib import Path


AGENT_TEMPLATE = """---
name: {agent_name}
description: "[TODO: Specific domain expertise description. Include 'Use PROACTIVELY when...' triggers. Under 200 characters recommended.]"
tools: [TODO: List tools - e.g., Read, Write, Edit, Bash, Glob, Grep]
---

You are [TODO: an expert ROLE specializing in DOMAIN].

## Core Workflow

When invoked:

1. [TODO: First step - typically analysis/understanding]
2. [TODO: Second step - typically planning/design]
3. [TODO: Third step - typically implementation]
4. [TODO: Fourth step - typically verification]

## [TODO: Domain Area 1]

- [TODO: Specific guidance for this area]
- [TODO: Decision criteria or patterns]

## [TODO: Domain Area 2]

- [TODO: Approaches and best practices]

## Tool Selection

Essential tools:

- **[TODO: Tool]**: [Purpose]
- **[TODO: Tool]**: [Purpose]

Collaboration:

- **[TODO: other-agent]**: [When to delegate]

## Output Format

[TODO: Template or example of expected output structure]

## Key Principles

- [TODO: Core principle 1]
- [TODO: Core principle 2]
- [TODO: Core principle 3]

[TODO: Closing statement - a key principle or focus directive]
"""


def validate_agent_name(name):
    """Validate agent name format."""
    if not re.match(r'^[a-z][a-z0-9-]*[a-z0-9]$', name):
        return False, "Name must be kebab-case (lowercase letters, digits, hyphens)"
    if len(name) > 50:
        return False, "Name must be 50 characters or fewer"
    if '--' in name:
        return False, "Name must not contain consecutive hyphens"
    return True, ""


def init_agent(agent_name, path):
    """
    Initialize a new agent markdown file from template.

    Args:
        agent_name: Name of the agent (kebab-case)
        path: Directory where the agent file should be created

    Returns:
        Path to created agent file, or None if error
    """
    valid, error = validate_agent_name(agent_name)
    if not valid:
        print(f"Error: {error}")
        return None

    output_dir = Path(path).resolve()
    if not output_dir.exists():
        print(f"Error: Directory does not exist: {output_dir}")
        return None

    agent_file = output_dir / f"{agent_name}.md"
    if agent_file.exists():
        print(f"Error: Agent file already exists: {agent_file}")
        return None

    content = AGENT_TEMPLATE.format(agent_name=agent_name)

    try:
        agent_file.write_text(content)
        print(f"Created: {agent_file}")
    except Exception as e:
        print(f"Error creating file: {e}")
        return None

    print(f"\nAgent '{agent_name}' initialized at {agent_file}")
    print("\nNext steps:")
    print("1. Replace all [TODO: ...] placeholders")
    print("2. Write a specific description for auto-triggering")
    print("3. List only the tools the agent needs")
    print("4. Run validate_agent.py to check the result")

    return agent_file


def main():
    if len(sys.argv) < 4 or sys.argv[2] != '--path':
        print("Usage: init_agent.py <agent-name> --path <output-directory>")
        print("\nExamples:")
        print("  init_agent.py my-specialist --path .claude/agents")
        print("  init_agent.py data-analyst --path .claude.ko/agents")
        sys.exit(1)

    agent_name = sys.argv[1]
    path = sys.argv[3]

    result = init_agent(agent_name, path)
    sys.exit(0 if result else 1)


if __name__ == "__main__":
    main()
