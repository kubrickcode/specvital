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
