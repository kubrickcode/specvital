-- Create "user_analysis_history" table
CREATE TABLE "public"."user_analysis_history" (
  "user_id" uuid NOT NULL,
  "analysis_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "analysis_id"),
  CONSTRAINT "fk_user_analysis_history_analysis" FOREIGN KEY ("analysis_id") REFERENCES "public"."analyses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_analysis_history_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_analysis_history_analysis" to table: "user_analysis_history"
CREATE INDEX "idx_user_analysis_history_analysis" ON "public"."user_analysis_history" ("analysis_id");
-- Create index "idx_user_analysis_history_user" to table: "user_analysis_history"
CREATE INDEX "idx_user_analysis_history_user" ON "public"."user_analysis_history" ("user_id", "updated_at");
