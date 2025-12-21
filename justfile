set dotenv-load := true

root_dir := justfile_directory()

bootstrap: install-psql install-sqlc install-playwright install-oapi-codegen

deps: deps-root deps-frontend

deps-root:
    pnpm install

deps-frontend:
    cd src/frontend && pnpm install

dump-schema:
    PGPASSWORD=postgres pg_dump -h specvital-postgres -U postgres -d specvital --schema-only --no-owner --no-privileges -n public | grep -v '^\\\|^SET \|^SELECT ' > src/backend/internal/db/schema.sql

gen-sqlc:
    cd src/backend && sqlc generate

gen-api:
    cd src/backend && oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
    cd src/frontend && npx openapi-typescript ../backend/api/openapi.yaml -o lib/api/generated-types.ts

install-playwright:
    npx playwright install --with-deps chromium

install-psql:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v psql &> /dev/null; then
      DEBIAN_FRONTEND=noninteractive apt-get update && \
        apt-get -y install lsb-release wget && \
        wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
        echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list && \
        apt-get update && \
        apt-get -y install postgresql-client-16
    fi

install-sqlc:
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0

install-oapi-codegen:
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

kill-port port:
    #!/usr/bin/env bash
    set -euo pipefail
    pid=$(ss -tlnp | grep ":{{ port }} " | sed -n 's/.*pid=\([0-9]*\).*/\1/p' | head -1)
    if [ -n "$pid" ]; then
        echo "Killing process $pid on port {{ port }}"
        kill -9 $pid
    else
        echo "No process found on port {{ port }}"
    fi

lint target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just lint justfile
        just lint config
        just lint backend
        just lint frontend
        ;;
      justfile)
        just --fmt --unstable
        ;;
      config)
        npx prettier --write "**/*.{json,yml,yaml,md}"
        ;;
      backend)
        gofmt -w src/backend
        ;;
      frontend)
        npx prettier --write "src/frontend/**/*.{ts,tsx}"
        cd src/frontend && pnpm eslint --fix --max-warnings=0 .
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

# Local db migration always initializes the database
migrate-local:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Resetting local database..."
    PGPASSWORD=postgres psql -h local-postgres -U postgres -c "DROP DATABASE IF EXISTS specvital;"
    PGPASSWORD=postgres psql -h local-postgres -U postgres -c "CREATE DATABASE specvital;"
    echo "Applying schema..."
    PGPASSWORD=postgres psql -h local-postgres -U postgres -d specvital -f src/backend/internal/db/schema.sql
    echo "✅ Migration complete!"

run target *args:
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      backend)
        cd src/backend
        if [[ " {{ args }} " =~ " --integration " ]]; then
          air
        else
          DATABASE_URL="$LOCAL_DATABASE_URL" \
          air
        fi
        ;;
      frontend)
        cd src/frontend && pnpm dev
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

sync-docs:
    baedal specvital/specvital.github.io/docs docs --exclude ".vitepress/**"

build target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just build backend
        just build frontend
        ;;
      backend)
        cd src/backend && go build ./...
        ;;
      frontend)
        cd src/frontend && pnpm build
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

release:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "⚠️  WARNING: This will trigger a production release!"
    echo ""
    echo "GitHub Actions will automatically:"
    echo "  - Analyze commits to determine version bump"
    echo "  - Generate release notes"
    echo "  - Create tag and GitHub release"
    echo "  - Update CHANGELOG.md"
    echo ""
    echo "Progress: https://github.com/specvital/web/actions"
    echo ""
    read -p "Type 'yes' to continue: " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Aborted."
        exit 1
    fi
    git checkout release
    git merge main
    git push origin release
    git checkout main
    echo "✅ Release triggered! Check GitHub Actions for progress."

run-worker:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -d "/tmp/collector" ]; then
        echo "Cloning collector repository..."
        git clone https://github.com/specvital/collector.git /tmp/collector
    else
        echo "Updating collector repository..."
        cd /tmp/collector && git pull
    fi
    cd /tmp/collector/src && \
    DATABASE_URL="$LOCAL_DATABASE_URL" \
    GITHUB_TOKEN="$GH_TOKEN" \
    go run ./cmd/worker

test target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just test backend
        just test frontend
        ;;
      backend)
        cd src/backend && go test -v ./...
        ;;
      frontend)
        cd src/frontend && pnpm test
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

update-core:
    cd src/backend && GOPROXY=direct go get -u github.com/specvital/core@main && go mod tidy
