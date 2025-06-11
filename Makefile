migrate:
	@echo "Migrating database..."
	@sql-migrate up -config="config.yaml" -env="db" \
		$(if $(version), -version $(version)) \
		$(if $(limit), -limit $(limit)) \
		$(if $(dryrun), -dryrun)

migrate-status:
	@sql-migrate status -config="config.yaml" -env="db"

migrate-down:
	@echo "Dropping database..."
	@sql-migrate down -config="config.yaml" -env="db" \
		$(if $(version), -version $(version)) \
		$(if $(limit), -limit $(limit)) \
		$(if $(dryrun), -dryrun)

migrate-new:
	@sql-migrate new -config="config.yaml" -env="db" $(name)

migrate-redo:
	@echo "Reapply last migration..."
	@sql-migrate redo -config="config.yaml" -env="db" \
		$(if $(dryrun), -dryrun)

bob:
	@echo "Running bob psql model codegen..."
	@go run github.com/stephenafamo/bob/gen/bobgen-psql@latest