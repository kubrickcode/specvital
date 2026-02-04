set dotenv-load := true

root_dir := justfile_directory()

deps: deps-root install-tbls

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
        npx prettier --cache --write "**/*.{json,yml,yaml,md}"
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
      */justfile|justfile)
        just --fmt --unstable
        ;;
      *.json|*.yml|*.yaml|*.md)
        npx prettier --cache --write "{{ file }}"
        ;;
      *)
        echo "No lint rule for: {{ file }}"
        ;;
    esac

makemigration name="changes":
    #!/usr/bin/env bash
    set -euo pipefail
    cd db
    atlas migrate diff {{ name }} --env local

    # Find the latest migration file
    LATEST=$(ls -t schema/migrations/*.sql | head -1)

    # Remove River-related DROP statements
    if grep -q "river_" "$LATEST"; then
        sed -i '/^-- Drop.*river/d' "$LATEST"
        sed -i '/^DROP TABLE.*river/d' "$LATEST"
        sed -i '/^DROP TYPE.*river/d' "$LATEST"
        sed -i '/^DROP FUNCTION.*river/d' "$LATEST"
        # Remove empty lines at end of file
        sed -i -e :a -e '/^\s*$/d;N;ba' "$LATEST"
        # Recalculate hash
        atlas migrate hash --env local
        echo "✅ Removed River DROP statements from $LATEST"
    fi

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

reset:
    #!/usr/bin/env bash
    set -euo pipefail
    psql "$DATABASE_URL" -c "DROP SCHEMA IF EXISTS atlas_schema_revisions CASCADE; DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
    cd db && atlas migrate apply --env local --allow-dirty

sync-docs:
    baedal specvital/specvital.github.io/docs docs --exclude ".vitepress/**"

erd:
    cd db && atlas schema inspect --env local --web

install-tbls:
    go install github.com/k1LoW/tbls@v1.92.3

gen-schema-docs:
    cd db && tbls doc --force

river_version := "v0.26.0"

river-install:
    go install github.com/riverqueue/river/cmd/river@{{ river_version }}

river-list:
    river migrate-get --line main --all --up | grep -E "^-- River main migration" || echo "Run 'just river-install' first"

river-dump:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Extracting River migration SQL ({{ river_version }})..."
    river migrate-get --line main --all --exclude-version 1 --up > /tmp/river_up.sql
    river migrate-get --line main --all --exclude-version 1 --down > /tmp/river_down.sql
    echo "✅ Exported to /tmp/river_up.sql and /tmp/river_down.sql"
    echo ""
    echo "Next steps:"
    echo "  1. Review the SQL files"
    echo "  2. Run: just river-migrate <name>"

river-migrate name="add_river_tables":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f /tmp/river_up.sql ]; then
        echo "❌ Run 'just river-dump' first"
        exit 1
    fi

    TIMESTAMP=$(date +%Y%m%d%H%M%S)
    TIMESTAMP2=$(printf "%d" $((TIMESTAMP + 1)))
    mkdir -p db/schema/rollbacks

    UP_FILE1="db/schema/migrations/${TIMESTAMP}_{{ name }}_1.sql"
    UP_FILE2="db/schema/migrations/${TIMESTAMP2}_{{ name }}_2.sql"
    DOWN_FILE="db/schema/rollbacks/${TIMESTAMP}_{{ name }}.down.sql"

    # Split SQL: Part 1 ends before river_job_state_in_bitmask function
    # Part 2 starts from that function
    SPLIT_MARKER="CREATE OR REPLACE FUNCTION river_job_state_in_bitmask"

    {
        echo "-- River Job Queue Tables (Part 1: Schema)"
        echo "-- River Version: {{ river_version }}"
        echo "-- https://github.com/riverqueue/river"
        echo ""
        sed 's|/\* TEMPLATE: schema \*/||g' /tmp/river_up.sql | awk -v marker="$SPLIT_MARKER" '$0 ~ marker {exit} {print}'
    } > "$UP_FILE1"

    {
        echo "-- River Job Queue Tables (Part 2: Functions)"
        echo "-- River Version: {{ river_version }}"
        echo "-- https://github.com/riverqueue/river"
        echo ""
        sed 's|/\* TEMPLATE: schema \*/||g' /tmp/river_up.sql | awk -v marker="$SPLIT_MARKER" 'found; $0 ~ marker {found=1; print}'
    } > "$UP_FILE2"

    {
        echo "-- River Job Queue Tables (Down)"
        echo "-- River Version: {{ river_version }}"
        echo ""
        sed 's|/\* TEMPLATE: schema \*/||g' /tmp/river_down.sql
    } > "$DOWN_FILE"

    cd db && atlas migrate hash --env local

    echo "✅ Created migration files:"
    echo "   $UP_FILE1"
    echo "   $UP_FILE2"
    echo "   $DOWN_FILE"
    echo ""
    echo "Next: just migrate"
