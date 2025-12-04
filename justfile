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
        just lint go
        ;;
      justfile)
        just --fmt --unstable
        ;;
      config)
        npx prettier --write "**/*.{json,yml,yaml,md}"
        ;;
      go)
        cd {{ root_dir }} && gofmt -w ./pkg
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

test:
    cd {{ root_dir }} && go test ./...

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
    echo "Progress: https://github.com/specvital/core/actions"
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
