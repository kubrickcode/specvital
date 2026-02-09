-- Create "user_bookmarks" table
CREATE TABLE "public"."user_bookmarks" (
  "user_id" uuid NOT NULL,
  "codebase_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "codebase_id"),
  CONSTRAINT "fk_user_bookmarks_codebase" FOREIGN KEY ("codebase_id") REFERENCES "public"."codebases" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_bookmarks_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_bookmarks_user" to table: "user_bookmarks"
CREATE INDEX "idx_user_bookmarks_user" ON "public"."user_bookmarks" ("user_id", "created_at");
