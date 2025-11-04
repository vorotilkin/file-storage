schema "public" {}

table "files" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "bucket" {
    type = text
    null = false
  }

  column "object_key" {
    type = text
    null = false
  }

  column "filename" {
    type = text
    null = false
  }

  column "content_type" {
    type = text
    null = false
  }

  column "created_at" {
    type     = timestamptz
    null     = false
    default = sql("CURRENT_TIMESTAMP")
  }

  column "uploaded_at" {
    type     = timestamptz
    null     = true
  }

  primary_key {
    columns = [column.id]
  }

  index "uq_files_bucket_object_key" {
    unique  = true
    columns = [column.bucket, column.object_key]
  }
}