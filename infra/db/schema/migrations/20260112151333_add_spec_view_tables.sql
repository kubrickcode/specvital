-- Create "spec_documents" table
CREATE TABLE "public"."spec_documents" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "analysis_id" uuid NOT NULL,
  "content_hash" bytea NOT NULL,
  "language" character varying(10) NOT NULL DEFAULT 'en',
  "executive_summary" text NULL,
  "model_id" character varying(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_spec_documents_hash_lang_model" UNIQUE ("content_hash", "language", "model_id"),
  CONSTRAINT "fk_spec_documents_analysis" FOREIGN KEY ("analysis_id") REFERENCES "public"."analyses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_spec_documents_analysis" to table: "spec_documents"
CREATE INDEX "idx_spec_documents_analysis" ON "public"."spec_documents" ("analysis_id");
-- Create "spec_domains" table
CREATE TABLE "public"."spec_domains" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "document_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  "classification_confidence" numeric(3,2) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_spec_domains_document" FOREIGN KEY ("document_id") REFERENCES "public"."spec_documents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_spec_domains_document_sort" to table: "spec_domains"
CREATE INDEX "idx_spec_domains_document_sort" ON "public"."spec_domains" ("document_id", "sort_order");
-- Create "spec_features" table
CREATE TABLE "public"."spec_features" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "domain_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_spec_features_domain" FOREIGN KEY ("domain_id") REFERENCES "public"."spec_domains" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_spec_features_domain_sort" to table: "spec_features"
CREATE INDEX "idx_spec_features_domain_sort" ON "public"."spec_features" ("domain_id", "sort_order");
-- Create "spec_behaviors" table
CREATE TABLE "public"."spec_behaviors" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "feature_id" uuid NOT NULL,
  "source_test_case_id" uuid NULL,
  "original_name" character varying(2000) NOT NULL,
  "converted_description" text NOT NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_spec_behaviors_feature" FOREIGN KEY ("feature_id") REFERENCES "public"."spec_features" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_spec_behaviors_test_case" FOREIGN KEY ("source_test_case_id") REFERENCES "public"."test_cases" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_spec_behaviors_feature_sort" to table: "spec_behaviors"
CREATE INDEX "idx_spec_behaviors_feature_sort" ON "public"."spec_behaviors" ("feature_id", "sort_order");
-- Create index "idx_spec_behaviors_source" to table: "spec_behaviors"
CREATE INDEX "idx_spec_behaviors_source" ON "public"."spec_behaviors" ("source_test_case_id") WHERE (source_test_case_id IS NOT NULL);
-- Drop "spec_view_cache" table
DROP TABLE "public"."spec_view_cache";
