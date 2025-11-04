-- Create "files" table
CREATE TABLE "files" ("id" serial NOT NULL, "bucket" text NOT NULL, "object_key" text NOT NULL, "filename" text NOT NULL, "content_type" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "uploaded_at" timestamptz NULL, PRIMARY KEY ("id"));
-- Create index "uq_files_bucket_object_key" to table: "files"
CREATE UNIQUE INDEX "uq_files_bucket_object_key" ON "files" ("bucket", "object_key");
