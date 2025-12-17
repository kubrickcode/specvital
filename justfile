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
        npx prettier --write "**/*.{json,yml,yaml,md}"
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

makemigration name="changes":
    cd db && atlas migrate diff {{ name }} --env local

migrate:
    cd db && atlas migrate apply --env local --allow-dirty

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
    echo "Progress: https://github.com/specvital/infra/actions"
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

reset target="all":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{ target }}" in
      all)
        just reset db
        just reset redis
        ;;
      db)
        psql "$DATABASE_URL" -c "DROP SCHEMA IF EXISTS atlas_schema_revisions CASCADE; DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
        cd db && atlas migrate apply --env local --allow-dirty
        ;;
      redis)
        redis-cli -u "$REDIS_URL" FLUSHALL
        ;;
      *)
        echo "Unknown target: {{ target }}"
        exit 1
        ;;
    esac

sync-docs:
    baedal specvital/specvital.github.io/docs docs --exclude ".vitepress/**"
