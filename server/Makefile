POSTGRES_MIGARTE_DB=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable # for docker-compose only

postgres_up:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) up

postgres_down:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) down

gen_ssl: gen_ssl_auth gen_ssl_notifications

gen_ssl_auth:
	cd cert &&  ./gen_auth_service.sh; cd .. && \
 	mkdir -p ./services/auth_service/configs/cert && \
	cp cert/server-auth-cert.pem cert/server-auth-key.pem cert/ca-auth-cert.pem ./services/auth_service/configs/cert && \
	mkdir -p ./services/chat_service/configs/cert && \
	cp cert/ca-auth-cert.pem cert/client-auth-cert.pem cert/client-auth-key.pem ./services/chat_service/configs/cert

gen_ssl_notifications:
	cd cert &&  ./gen_notifications_service.sh; cd .. && \
 	mkdir -p ./services/notifications_service/configs/cert && \
	cp cert/server-notifications-cert.pem cert/server-notifications-key.pem cert/ca-notifications-cert.pem ./services/notifications_service/configs/cert && \
	mkdir -p ./services/auth_service/configs/cert && \
	cp cert/ca-notifications-cert.pem cert/client-notifications-cert.pem cert/client-notifications-key.pem ./services/auth_service/configs/cert