# public.codebases

## Description

## Columns

| Name             | Type                     | Default                         | Nullable | Children                                                                                | Parents | Comment |
| ---------------- | ------------------------ | ------------------------------- | -------- | --------------------------------------------------------------------------------------- | ------- | ------- |
| id               | uuid                     | gen_random_uuid()               | false    | [public.analyses](public.analyses.md) [public.user_bookmarks](public.user_bookmarks.md) |         |         |
| host             | varchar(255)             | 'github.com'::character varying | false    |                                                                                         |         |         |
| owner            | varchar(255)             |                                 | false    |                                                                                         |         |         |
| name             | varchar(255)             |                                 | false    |                                                                                         |         |         |
| default_branch   | varchar(100)             |                                 | true     |                                                                                         |         |         |
| created_at       | timestamp with time zone | now()                           | false    |                                                                                         |         |         |
| updated_at       | timestamp with time zone | now()                           | false    |                                                                                         |         |         |
| last_viewed_at   | timestamp with time zone |                                 | true     |                                                                                         |         |         |
| external_repo_id | varchar(64)              |                                 | false    |                                                                                         |         |         |
| is_stale         | boolean                  | false                           | false    |                                                                                         |         |         |
| is_private       | boolean                  | false                           | false    |                                                                                         |         |         |

## Constraints

| Name           | Type        | Definition       |
| -------------- | ----------- | ---------------- |
| codebases_pkey | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name                           | Definition                                                                                                                 |
| ------------------------------ | -------------------------------------------------------------------------------------------------------------------------- |
| codebases_pkey                 | CREATE UNIQUE INDEX codebases_pkey ON public.codebases USING btree (id)                                                    |
| idx_codebases_owner_name       | CREATE INDEX idx_codebases_owner_name ON public.codebases USING btree (owner, name)                                        |
| idx_codebases_last_viewed      | CREATE INDEX idx_codebases_last_viewed ON public.codebases USING btree (last_viewed_at) WHERE (last_viewed_at IS NOT NULL) |
| idx_codebases_external_repo_id | CREATE UNIQUE INDEX idx_codebases_external_repo_id ON public.codebases USING btree (host, external_repo_id)                |
| idx_codebases_identity         | CREATE UNIQUE INDEX idx_codebases_identity ON public.codebases USING btree (host, owner, name) WHERE (is_stale = false)    |
| idx_codebases_public           | CREATE INDEX idx_codebases_public ON public.codebases USING btree (is_private) WHERE (is_private = false)                  |

## Relations

```mermaid
erDiagram

"public.analyses" }o--|| "public.codebases" : "FOREIGN KEY (codebase_id) REFERENCES codebases(id) ON DELETE CASCADE"
"public.user_bookmarks" }o--|| "public.codebases" : "FOREIGN KEY (codebase_id) REFERENCES codebases(id) ON DELETE CASCADE"

"public.codebases" {
  uuid id
  varchar_255_ host
  varchar_255_ owner
  varchar_255_ name
  varchar_100_ default_branch
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone last_viewed_at
  varchar_64_ external_repo_id
  boolean is_stale
  boolean is_private
}
"public.analyses" {
  uuid id
  uuid codebase_id FK
  varchar_40_ commit_sha
  varchar_255_ branch_name
  analysis_status status
  text error_message
  timestamp_with_time_zone started_at
  timestamp_with_time_zone completed_at
  timestamp_with_time_zone created_at
  integer total_suites
  integer total_tests
  timestamp_with_time_zone committed_at
  varchar_100_ parser_version
}
"public.user_bookmarks" {
  uuid user_id FK
  uuid codebase_id FK
  timestamp_with_time_zone created_at
  uuid id
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
