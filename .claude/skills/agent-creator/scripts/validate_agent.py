#!/usr/bin/env python3
"""
Agent Validator - Validates agent markdown file structure

Usage:
    validate_agent.py <path/to/agent.md>
    validate_agent.py <path/to/agents-directory>

Examples:
    validate_agent.py .claude/agents/my-agent.md
    validate_agent.py .claude/agents/
"""

import sys
import re
from pathlib import Path

VALID_TOOLS = {
    "Read", "Write", "Edit", "Bash", "Glob", "Grep",
    "WebFetch", "WebSearch", "AskUserQuestion", "TodoWrite",
    "NotebookEdit", "Task",
}

VALID_MODELS = {"haiku", "sonnet", "opus", "inherit"}

VALID_PERMISSION_MODES = {
    "default", "acceptEdits", "delegate", "dontAsk",
    "bypassPermissions", "plan",
}


def parse_frontmatter(content):
    """Extract YAML frontmatter from markdown content."""
    match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    if not match:
        return None, "No valid YAML frontmatter found (must start with ---)"

    frontmatter_text = match.group(1)
    fields = {}
    for line in frontmatter_text.strip().split('\n'):
        line = line.strip()
        if not line or line.startswith('#'):
            continue
        if ':' in line:
            key, _, value = line.partition(':')
            fields[key.strip()] = value.strip().strip('"').strip("'")

    return fields, None


def validate_single(filepath):
    """Validate a single agent markdown file. Returns list of (level, message)."""
    issues = []
    path = Path(filepath)

    if not path.exists():
        return [("ERROR", f"File not found: {filepath}")]
    if not path.suffix == '.md':
        return [("ERROR", f"Agent file must be .md: {filepath}")]

    content = path.read_text()

    # Parse frontmatter
    fields, error = parse_frontmatter(content)
    if error:
        return [("ERROR", error)]

    # Required: name
    if 'name' not in fields:
        issues.append(("ERROR", "Missing required field: name"))
    else:
        name = fields['name']
        expected_name = path.stem
        if name != expected_name:
            issues.append(("ERROR", f"name '{name}' does not match filename '{expected_name}'"))
        if not re.match(r'^[a-z][a-z0-9-]*[a-z0-9]$', name):
            issues.append(("ERROR", f"name '{name}' is not valid kebab-case"))

    # Required: description
    if 'description' not in fields:
        issues.append(("ERROR", "Missing required field: description"))
    else:
        desc = fields['description']
        if '[TODO' in desc:
            issues.append(("ERROR", "description contains [TODO] placeholder"))
        if len(desc) < 20:
            issues.append(("WARNING", f"description is very short ({len(desc)} chars). Aim for 50-200 chars."))
        if len(desc) > 300:
            issues.append(("WARNING", f"description is long ({len(desc)} chars). Consider trimming to under 200."))
        if 'proactively' not in desc.lower() and 'use when' not in desc.lower() and 'use for' not in desc.lower():
            issues.append(("WARNING", "description lacks trigger language (e.g., 'Use PROACTIVELY when...')"))

    # Optional: tools
    if 'tools' in fields:
        tools_str = fields['tools']
        if '[TODO' not in tools_str:
            tools = [t.strip() for t in tools_str.split(',')]
            for tool in tools:
                # Allow Task(agent-name) syntax
                tool_name = re.match(r'^(\w+)', tool)
                if tool_name and tool_name.group(1) not in VALID_TOOLS:
                    issues.append(("WARNING", f"Unknown tool: {tool}"))
    else:
        issues.append(("WARNING", "No tools field — agent will inherit ALL tools. Consider explicit scoping."))

    # Optional: model
    if 'model' in fields:
        if fields['model'] not in VALID_MODELS:
            issues.append(("WARNING", f"Unknown model: {fields['model']}. Valid: {', '.join(sorted(VALID_MODELS))}"))

    # Optional: permissionMode
    if 'permissionMode' in fields:
        if fields['permissionMode'] not in VALID_PERMISSION_MODES:
            issues.append(("WARNING", f"Unknown permissionMode: {fields['permissionMode']}"))

    # Body checks
    body_match = re.match(r'^---\s*\n.*?\n---\s*\n(.*)$', content, re.DOTALL)
    if body_match:
        body = body_match.group(1).strip()

        if not body:
            issues.append(("ERROR", "Agent body (system prompt) is empty"))
        elif '[TODO' in body:
            todo_count = body.count('[TODO')
            issues.append(("WARNING", f"Body contains {todo_count} [TODO] placeholder(s)"))

        lines = body.split('\n')
        line_count = len(lines)

        if line_count < 5:
            issues.append(("WARNING", f"Body is very short ({line_count} lines). Consider adding more guidance."))

        # Check for opening role statement
        first_line = body.split('\n')[0].strip() if body else ""
        if first_line and not first_line.lower().startswith('you are'):
            issues.append(("SUGGESTION", "Body should start with 'You are a [ROLE]...' for clear expertise establishment"))

        # Check for H4 usage (discouraged)
        if re.search(r'^####\s', body, re.MULTILINE):
            issues.append(("SUGGESTION", "Avoid H4 (####) headers. Use H2/H3 only for better scannability."))
    else:
        issues.append(("ERROR", "Could not parse agent body after frontmatter"))

    return issues


def main():
    if len(sys.argv) < 2:
        print("Usage: validate_agent.py <path/to/agent.md or agents-directory>")
        sys.exit(1)

    target = Path(sys.argv[1])
    files = []

    if target.is_dir():
        files = sorted(target.glob("*.md"))
        if not files:
            print(f"No .md files found in {target}")
            sys.exit(1)
    elif target.is_file():
        files = [target]
    else:
        print(f"Path not found: {target}")
        sys.exit(1)

    total_errors = 0
    total_warnings = 0
    total_suggestions = 0

    for f in files:
        issues = validate_single(f)
        if not issues:
            print(f"PASS: {f.name}")
            continue

        has_error = any(level == "ERROR" for level, _ in issues)
        print(f"\n{'FAIL' if has_error else 'WARN'}: {f.name}")
        for level, msg in issues:
            prefix = {"ERROR": "  [x]", "WARNING": "  [!]", "SUGGESTION": "  [~]"}[level]
            print(f"{prefix} {msg}")
            if level == "ERROR":
                total_errors += 1
            elif level == "WARNING":
                total_warnings += 1
            else:
                total_suggestions += 1

    print(f"\n--- Summary ---")
    print(f"Files: {len(files)}")
    print(f"Errors: {total_errors}, Warnings: {total_warnings}, Suggestions: {total_suggestions}")

    sys.exit(1 if total_errors > 0 else 0)


if __name__ == "__main__":
    main()
