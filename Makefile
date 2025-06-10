migrate:
	@echo "Migrating database..."
	@sql-migrate up -config="config.yaml" -env="db"

migrate\:status:
	@sql-migrate status -config="config.yaml" -env="db"

migrate\:down:
	@sql-migrate down -config="config.yaml" -env="db"

migrate\:new:
	@sql-migrate new -config="config.yaml" -env="db" -seq $(name)

migrate:\redo:
	@sql-migrate redo -config="config.yaml" -env="db"