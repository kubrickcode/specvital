#!/bin/bash

cat > /dev/null

MESSAGE="${1:-✅ Work completed!}"
PROJECT_NAME=$(basename "$PWD")
FULL_MESSAGE="[$PROJECT_NAME] $MESSAGE"

curl -s -X POST \
  -H 'Content-type: application/json' \
  --data "{\"content\":\"$FULL_MESSAGE\"}" \
  "$DISCORD_NOTIFY_WEBHOOK_URL" || true
