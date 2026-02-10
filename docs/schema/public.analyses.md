# public.analyses

## Description

## Columns

| Name           | Type                     | Default                     | Nullable | Children                                                                                                                                                                                                  | Parents                                 | Comment |
| -------------- | ------------------------ | --------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------- | ------- |
| id             | uuid                     | gen_random_uuid()           | false    | [public.user_analysis_history](public.user_analysis_history.md) [public.test_files](public.test_files.md) [public.spec_documents](public.spec_documents.md) [public.usage_events](public.usage_events.md) |                                         |         |
| codebase_id    | uuid                     |                             | false    |                                                                                                                                                                                                           | [public.codebases](public.codebases.md) |         |
| commit_sha     | varchar(40)              |                             | false    |                                                                                                                                                                                                           |                                         |         |
| branch_name    | varchar(255)             |                             | true     |                                                                                                                                                                                                           |                                         |         |
| status         | analysis_status          | 'pending'::analysis_status  | false    |                                                                                                                                                                                                           |                                         |         |
| error_message  | text                     |                             | true     |                                                                                                                                                                                                           |                                         |         |
| started_at     | timestamp with time zone |                             | true     |                                                                                                                                                                                                           |                                         |         |
| completed_at   | timestamp with time zone |                             | true     |                                                                                                                                                                                                           |                                         |         |
| created_at     | timestamp with time zone | now()                       | false    |                                                                                                                                                                                                           |                                         |         |
| total_suites   | integer                  | 0                           | false    |                                                                                                                                                                                                           |                                         |         |
| total_tests    | integer                  | 0                           | false    |                                                                                                                                                                                                           |                                         |         |
| committed_at   | timestamp with time zone |                             | true     |                                                                                                                                                                                                           |                                         |         |
| parser_version | varchar(100)             | 'legacy'::character varying | false    |                                                                                                                                                                                                           |                                         |         |

## Constraints

| Name                 | Type        | Definition                                                           |
| -------------------- | ----------- | -------------------------------------------------------------------- |
| fk_analyses_codebase | FOREIGN KEY | FOREIGN KEY (codebase_id) REFERENCES codebases(id) ON DELETE CASCADE |
| analyses_pkey        | PRIMARY KEY | PRIMARY KEY (id)                                                     |

## Indexes

| Name                                 | Definition                                                                                                                                                                      |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| analyses_pkey                        | CREATE UNIQUE INDEX analyses_pkey ON public.analyses USING btree (id)                                                                                                           |
| idx_analyses_codebase_status         | CREATE INDEX idx_analyses_codebase_status ON public.analyses USING btree (codebase_id, status)                                                                                  |
| idx_analyses_created                 | CREATE INDEX idx_analyses_created ON public.analyses USING btree (codebase_id, created_at)                                                                                      |
| uq_analyses_completed_commit_version | CREATE UNIQUE INDEX uq_analyses_completed_commit_version ON public.analyses USING btree (codebase_id, commit_sha, parser_version) WHERE (status = 'completed'::analysis_status) |

## Relations

```mermaid
erDiagram

"public.user_analysis_history" }o--|| "public.analyses" : "FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE CASCADE"
"public.test_files" }o--|| "public.analyses" : "FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE CASCADE"
"public.spec_documents" }o--|| "public.analyses" : "FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE CASCADE"
"public.usage_events" }o--o| "public.analyses" : "FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE SET NULL"
"public.analyses" }o--|| "public.codebases" : "FOREIGN KEY (codebase_id) REFERENCES codebases(id) ON DELETE CASCADE"

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
"public.user_analysis_history" {
  uuid user_id FK
  uuid analysis_id FK
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  uuid id
  integer retention_days_at_creation
}
"public.test_files" {
  uuid id
  uuid analysis_id FK
  varchar_1000_ file_path
  varchar_50_ framework
  jsonb domain_hints
}
"public.spec_documents" {
  uuid id
  uuid analysis_id FK
  bytea content_hash
  varchar_10_ language
  text executive_summary
  varchar_100_ model_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  integer version
  uuid user_id FK
  integer retention_days_at_creation
}
"public.usage_events" {
  uuid id
  uuid user_id FK
  usage_event_type event_type
  uuid analysis_id FK
  uuid document_id FK
  integer quota_amount
  timestamp_with_time_zone created_at
}
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
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
