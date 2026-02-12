set dotenv-load := true

root_dir := justfile_directory()

bootstrap: install-atlas install-psql install-sqlc install-docker install-playwright install-oapi-codegen install-tbls

deps: deps-root deps-web-frontend

migrate:
    #!/usr/bin/env bash
    set -euo pipefail
    cd {{ root_dir }}/infra && just migrate
    cd {{ root_dir }}/apps/web && just dump-schema
    cd {{ root_dir }}/apps/worker && just dump-schema
    echo "Migration complete! Schemas dumped to web and worker."

deps-root:
    pnpm install

deps-web-frontend:
    cd {{ root_dir }}/apps/web && just deps-frontend

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

install-docker:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v docker &> /dev/null; then
      curl -fsSL https://get.docker.com | sh
    fi

install-playwright:
    cd {{ root_dir }}/apps/web && npx playwright install --with-deps chromium

install-oapi-codegen:
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1

install-atlas:
    curl -sSf https://atlasgo.sh | sh

install-tbls:
    go install github.com/k1LoW/tbls@v1.92.3

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

_lint-justfile mode *target:
    just --fmt --unstable

_lint-config mode *target:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{ mode }}" = "batch" ]; then
        npx prettier --write --cache "**/*.{json,yml,yaml,md}"
    else
        npx prettier --write --cache "{{ target }}"
    fi

_lint-ts mode *target:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{ mode }}" = "batch" ]; then
        npx prettier --write --cache "**/*.{ts,tsx}"
        cd apps/web/frontend && npx eslint --fix --max-warnings=0 .
    else
        npx prettier --write --cache "{{ target }}"
        if [[ "{{ target }}" == apps/web/frontend/* ]]; then
            cd apps/web/frontend && npx eslint --fix --max-warnings=0 "{{ root_dir }}/{{ target }}"
        fi
    fi

_lint-go mode *target:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{ mode }}" = "batch" ]; then
        gofmt -w .
        for mod in apps/web/backend apps/worker lib; do
            (cd "$mod" && go vet ./...)
        done
    else
        gofmt -w "{{ target }}"
        abs_dir="$(cd "$(dirname '{{ target }}')" && pwd)"
        output=$( (cd "$abs_dir" && go vet .) 2>&1) || {
            if echo "$output" | grep -q "build constraints exclude all Go files"; then
                true
            else
                echo "$output" >&2
                exit 1
            fi
        }
    fi

lint target="all":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ target }}" in
      all)
        just _lint-justfile batch
        just _lint-config batch
        just _lint-go batch
        just _lint-ts batch
        ;;
      justfile)  just _lint-justfile batch ;;
      config)    just _lint-config batch ;;
      go)              just _lint-go batch ;;
      ts|web-frontend) just _lint-ts batch ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

lint-file file:
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ file }}" in
      */justfile|*Justfile)      just _lint-justfile file "{{ file }}" ;;
      *.json|*.yml|*.yaml|*.md)  just _lint-config file "{{ file }}" ;;
      *.ts|*.tsx)                just _lint-ts file "{{ file }}" ;;
      *.go)                      just _lint-go file "{{ file }}" ;;
      *) ;;
    esac

release:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "This will trigger a production release!"
    echo "GitHub Actions: https://github.com/kubrickcode/specvital/actions"
    read -p "Type 'yes' to continue: " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Aborted."
        exit 1
    fi
    git checkout release
    git merge main
    git push origin release
    git checkout main
    echo "Release triggered! Check GitHub Actions for progress."

run target *args:
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ target }}" in
      web-backend)
        cd {{ root_dir }}/apps/web && just run backend {{ args }}
        ;;
      web-frontend)
        cd {{ root_dir }}/apps/web && just run frontend {{ args }}
        ;;
      analyzer)
        cd {{ root_dir }}/apps/worker && just run-analyzer {{ args }}
        ;;
      spec-generator)
        cd {{ root_dir }}/apps/worker && just run-spec-generator {{ args }}
        ;;
      spec-generator-mock)
        cd {{ root_dir }}/apps/worker && just run-spec-generator-mock {{ args }}
        ;;
      *)
        echo "Unknown: {{ target }}. Use: web-backend, web-frontend, analyzer, spec-generator, spec-generator-mock"
        exit 1
        ;;
    esac

test project="all" target="all":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ project }}" in
      web)
        cd {{ root_dir }}/apps/web && just test {{ target }}
        ;;
      worker)
        cd {{ root_dir }}/apps/worker && just test {{ target }}
        ;;
      core)
        cd {{ root_dir }}/lib && just test {{ target }}
        ;;
      all)
        cd {{ root_dir }}/apps/web && just test all
        cd {{ root_dir }}/apps/worker && just test all
        cd {{ root_dir }}/lib && just test all
        ;;
      *)
        echo "Unknown: {{ project }}. Use: web, worker, core, all"
        exit 1
        ;;
    esac
