#!/bin/bash
input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')
if echo "$command" | grep -qE '(git\s+(commit|push)|git\s+.*&&\s*git\s+(commit|push))'; then
  cat << 'EOF'
{"decision": "block", "reason": "Git commit/push requires explicit user request. Use '/commit' command or ask user for confirmation."}
EOF
  exit 0
fi
exit 0
