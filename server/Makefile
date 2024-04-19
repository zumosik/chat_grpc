POSTGRES_MIGARTE_DB=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable # for docker-compose only

postgres_up:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) up

postgres_down:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) down

gen_ssl:
	cd cert &&  ./gen.sh; cd .. && \
 	mkdir -p ./services/auth_service/configs/cert && \
	cp cert/server-auth-cert.pem cert/server-auth-key.pem cert/ca-auth-cert.pem ./services/auth_service/configs/cert && \
	mkdir -p ./services/chat_service/configs/cert && \
	cp cert/ca-auth-cert.pem cert/client-auth-cert.pem cert/client-auth-key.pem ./services/chat_service/configs/cert