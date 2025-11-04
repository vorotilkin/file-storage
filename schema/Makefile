DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=postgres
DB_PASSWORD?=password
DB_DATABASE?=file_storage
DB_SCHEMA?=public
DB_SSL_MODE?=disable
DB_DSN?=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?$(if $(DB_SSL_MODE),sslmode=$(DB_SSL_MODE),)
MIGRATION_NAME?=init
ATLASGO_BIN?=$(shell which atlas)
JET_BIN?=$(shell which jet)

atlas.apply:
	@$(ATLASGO_BIN) schema apply \
      -u "$(DB_DSN)" \
      --to file://schema.pg.hcl

atlas.inspect:
	@$(ATLASGO_BIN) \
    schema inspect \
      -u "$(DB_DSN)" > schema.pg.hcl

atlas.migration.make:
	@$(ATLASGO_BIN) migrate diff $(MIGRATION_NAME) \
      --dir "file://migrations" \
      --to "file://schema.pg.hcl" \
      --dev-url "docker://postgres/15/test?search_path=public"

jet.generate:
	@$(JET_BIN) -dsn=$(DB_DSN) -schema=$(DB_SCHEMA) -path=./gen