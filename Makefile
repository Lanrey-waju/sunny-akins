run:
	go run ./cmd/web

psql:
	psql ${SUNNY_DSN}

db/migrations/up:
	@echo 'Running up migrations'
	@goose -dir "sql/schema" up
	
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	goose -dir "sql/schema" -s create ${name} sql
