set dotenv-load := true

root_dir := justfile_directory()

deps: deps-root deps-frontend

deps-root:
    pnpm install

deps-frontend:
    cd src/frontend && pnpm install

lint target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just lint justfile
        just lint config
        just lint go
        ;;
      justfile)
        just --fmt --unstable
        ;;
      config)
        npx prettier --write "**/*.{json,yml,yaml,md}"
        ;;
      go)
        gofmt -w src/backend
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

run target:
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      backend)
        cd src/backend && air
        ;;
      frontend)
        cd src/frontend && pnpm dev
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

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
