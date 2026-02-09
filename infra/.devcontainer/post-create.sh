#!/bin/bash
set -euo pipefail

ATLAS_VERSION="0.38.0"

if ! command -v psql &> /dev/null || ! command -v redis-cli &> /dev/null; then
  DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get -y install lsb-release wget && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list && \
    apt-get update && \
    apt-get -y install postgresql-client-16 redis-tools
fi

curl -sSfL "https://release.ariga.io/atlas/atlas-linux-amd64-v${ATLAS_VERSION}" -o /usr/local/bin/atlas
chmod +x /usr/local/bin/atlas

[ -s /root/.claude.json ] || echo '{}' > /root/.claude.json

command -v claude &>/dev/null || curl -fsSL https://claude.ai/install.sh | bash

npm install -g baedal
