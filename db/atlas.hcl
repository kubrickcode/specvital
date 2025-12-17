env "local" {
  src = "file://schema/schema.hcl"
  url = "postgres://postgres:postgres@postgres:5432/specvital?sslmode=disable"
  dev = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"

  migration {
    dir = "file://schema/migrations"
    revisions_schema = "public"
  }
}

env "ci" {
  src = "file://schema/schema.hcl"
  url = "postgres://postgres:postgres@localhost:5432/specvital?sslmode=disable"
  dev = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

  migration {
    dir = "file://schema/migrations"
    revisions_schema = "public"
  }
}

env "production" {
  src = "file://schema/schema.hcl"
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://schema/migrations"
    # Explicitly configure revision table in public schema
    revisions_schema = "public"
  }

  # Enable migration versioning and safety checks
  diff {
    # Skip destructive changes detection for production
    skip {
      drop_schema = false
      drop_table  = false
    }
  }
}
