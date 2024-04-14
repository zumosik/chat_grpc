POSTGRES_MIGARTE_DB=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable # for docker-compose only

postgres_up:
	goose -dir migrations postgres $(POSTGRES_MIGARTE_DB) up

postgres_down:
	goose -dir migrations postgres $(POSTGRES_MIGARTE_DB) down