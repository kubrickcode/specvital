-- Create "user_specview_history" table
CREATE TABLE "public"."user_specview_history" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "document_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "uq_user_specview_history_user_document" UNIQUE ("user_id", "document_id"),
  CONSTRAINT "fk_user_specview_history_document" FOREIGN KEY ("document_id") REFERENCES "public"."spec_documents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_specview_history_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_specview_history_cursor" to table: "user_specview_history"
CREATE INDEX "idx_user_specview_history_cursor" ON "public"."user_specview_history" ("user_id", "updated_at", "id");
-- Create index "idx_user_specview_history_document" to table: "user_specview_history"
CREATE INDEX "idx_user_specview_history_document" ON "public"."user_specview_history" ("document_id");
