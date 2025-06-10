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
	@sql-migrate new -config="config.yaml" -env="db" -seq $(name)

migrate-redo:
	@sql-migrate redo -config="config.yaml" -env="db" \
		$(if $(dryrun), -dryrun)