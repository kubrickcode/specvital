# public.spec_documents

## Description

## Columns

| Name                       | Type                     | Default                 | Nullable | Children                                                                                                                                                    | Parents                               | Comment |
| -------------------------- | ------------------------ | ----------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------- | ------- |
| id                         | uuid                     | gen_random_uuid()       | false    | [public.spec_domains](public.spec_domains.md) [public.user_specview_history](public.user_specview_history.md) [public.usage_events](public.usage_events.md) |                                       |         |
| analysis_id                | uuid                     |                         | false    |                                                                                                                                                             | [public.analyses](public.analyses.md) |         |
| content_hash               | bytea                    |                         | false    |                                                                                                                                                             |                                       |         |
| language                   | varchar(10)              | 'en'::character varying | false    |                                                                                                                                                             |                                       |         |
| executive_summary          | text                     |                         | true     |                                                                                                                                                             |                                       |         |
| model_id                   | varchar(100)             |                         | false    |                                                                                                                                                             |                                       |         |
| created_at                 | timestamp with time zone | now()                   | false    |                                                                                                                                                             |                                       |         |
| updated_at                 | timestamp with time zone | now()                   | false    |                                                                                                                                                             |                                       |         |
| version                    | integer                  | 1                       | false    |                                                                                                                                                             |                                       |         |
| user_id                    | uuid                     |                         | false    |                                                                                                                                                             | [public.users](public.users.md)       |         |
| retention_days_at_creation | integer                  |                         | true     |                                                                                                                                                             |                                       |         |

## Constraints

| Name                                           | Type        | Definition                                                                         |
| ---------------------------------------------- | ----------- | ---------------------------------------------------------------------------------- |
| chk_retention_days_positive                    | CHECK       | CHECK (((retention_days_at_creation IS NULL) OR (retention_days_at_creation > 0))) |
| fk_spec_documents_analysis                     | FOREIGN KEY | FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE CASCADE                |
| fk_spec_documents_user                         | FOREIGN KEY | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE                       |
| spec_documents_pkey                            | PRIMARY KEY | PRIMARY KEY (id)                                                                   |
| uq_spec_documents_user_analysis_lang_version   | UNIQUE      | UNIQUE (user_id, analysis_id, language, version)                                   |
| uq_spec_documents_user_hash_lang_model_version | UNIQUE      | UNIQUE (user_id, content_hash, language, model_id, version)                        |

## Indexes

| Name                                           | Definition                                                                                                                                                   |
| ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| spec_documents_pkey                            | CREATE UNIQUE INDEX spec_documents_pkey ON public.spec_documents USING btree (id)                                                                            |
| idx_spec_documents_analysis                    | CREATE INDEX idx_spec_documents_analysis ON public.spec_documents USING btree (analysis_id)                                                                  |
| uq_spec_documents_user_analysis_lang_version   | CREATE UNIQUE INDEX uq_spec_documents_user_analysis_lang_version ON public.spec_documents USING btree (user_id, analysis_id, language, version)              |
| uq_spec_documents_user_hash_lang_model_version | CREATE UNIQUE INDEX uq_spec_documents_user_hash_lang_model_version ON public.spec_documents USING btree (user_id, content_hash, language, model_id, version) |
| idx_spec_documents_user_created                | CREATE INDEX idx_spec_documents_user_created ON public.spec_documents USING btree (user_id, created_at)                                                      |
| idx_spec_documents_content_hash_lang_model     | CREATE INDEX idx_spec_documents_content_hash_lang_model ON public.spec_documents USING btree (content_hash, language, model_id)                              |
| idx_spec_documents_retention_cleanup           | CREATE INDEX idx_spec_documents_retention_cleanup ON public.spec_documents USING btree (created_at) WHERE (retention_days_at_creation IS NOT NULL)           |

## Relations

```mermaid
erDiagram

"public.spec_domains" }o--|| "public.spec_documents" : "FOREIGN KEY (document_id) REFERENCES spec_documents(id) ON DELETE CASCADE"
"public.user_specview_history" }o--|| "public.spec_documents" : "FOREIGN KEY (document_id) REFERENCES spec_documents(id) ON DELETE CASCADE"
"public.usage_events" }o--o| "public.spec_documents" : "FOREIGN KEY (document_id) REFERENCES spec_documents(id) ON DELETE SET NULL"
"public.spec_documents" }o--|| "public.analyses" : "FOREIGN KEY (analysis_id) REFERENCES analyses(id) ON DELETE CASCADE"
"public.spec_documents" }o--|| "public.users" : "FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE"

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
"public.spec_domains" {
  uuid id
  uuid document_id FK
  varchar_255_ name
  text description
  integer sort_order
  numeric_3_2_ classification_confidence
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.user_specview_history" {
  uuid id
  uuid user_id FK
  uuid document_id FK
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
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
"public.users" {
  uuid id
  varchar_255_ email
  varchar_255_ username
  text avatar_url
  timestamp_with_time_zone last_login_at
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  integer token_version
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
