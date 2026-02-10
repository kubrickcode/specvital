set dotenv-load := true

root_dir := justfile_directory()

deps: deps-root

deps-root:
    pnpm install

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

lint-file file:
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ file }}" in
      */justfile|*Justfile)
        just --fmt --unstable
        ;;
      *.json|*.yml|*.yaml|*.md)
        npx prettier --write --cache "{{ file }}"
        ;;
      *.ts|*.tsx)
        npx prettier --write --cache "{{ file }}"
        ;;
      *.go)
        gofmt -w "{{ file }}"
        go vet "$(dirname '{{ file }}')/..."
        ;;
      *)
        ;;
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
      *)
        echo "Unknown: {{ target }}. Use: web-backend, web-frontend, analyzer, spec-generator"
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
        cd {{ root_dir }}/packages/core && just test {{ target }}
        ;;
      all)
        cd {{ root_dir }}/apps/web && just test all
        cd {{ root_dir }}/apps/worker && just test all
        cd {{ root_dir }}/packages/core && just test all
        ;;
      *)
        echo "Unknown: {{ project }}. Use: web, worker, core, all"
        exit 1
        ;;
    esac

lint target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just lint justfile
        just lint config
        just lint go
        just lint web-frontend
        ;;
      justfile)
        just --fmt --unstable
        ;;
      config)
        npx prettier --write --cache "**/*.{json,yml,yaml,md}"
        ;;
      go)
        gofmt -w .
        ;;
      web-frontend)
        cd apps/web && just lint-frontend
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac
