set dotenv-load := true

root_dir := justfile_directory()

deps: deps-root

deps-root:
    pnpm install

lint target="all":
    #!/usr/bin/env bash
    set -euox pipefail
    case "{{ target }}" in
      all)
        just lint justfile
        just lint config
        ;;
      justfile)
        just --fmt --unstable
        ;;
      config)
        npx prettier --write --cache "**/*.{json,yml,yaml,md}"
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

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
      *)
        echo "No lint rule for: {{ file }}"
        ;;
    esac
