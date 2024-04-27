POSTGRES_MIGARTE_DB=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable # for docker-compose only

postgres_up:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) up

postgres_down:
	goose -dir postgres_migrations postgres $(POSTGRES_MIGARTE_DB) down

gen_ssl: clear_ssl gen_ssl_auth gen_ssl_notifications gen_ssl_rooms gen_ssl_chat

clear_ssl:
	rm cert/*.pem

gen_ssl_rooms:
	cd cert &&  ./gen_service.sh rooms; cd .. && \
	mkdir -p ./services/rooms_service/configs/cert && \
	cp cert/server-rooms-cert.pem cert/server-rooms-key.pem cert/ca-rooms-cert.pem ./services/rooms_service/configs/cert && \
	cp cert/ca-rooms-cert.pem cert/client-rooms-cert.pem cert/client-rooms-key.pem ./services/chat_service/configs/cert

gen_ssl_auth:
	cd cert &&  ./gen_service.sh auth; cd .. && \
	cp cert/server-auth-cert.pem cert/server-auth-key.pem cert/ca-auth-cert.pem ./services/auth_service/configs/cert && \
	cp cert/ca-auth-cert.pem cert/client-auth-cert.pem cert/client-auth-key.pem ./services/chat_service/configs/cert && \
	cp cert/ca-auth-cert.pem cert/client-auth-cert.pem cert/client-auth-key.pem ./services/rooms_service/configs/cert


gen_ssl_notifications:
	cd cert &&  ./gen_service.sh notifications; cd .. && \
	cp cert/server-notifications-cert.pem cert/server-notifications-key.pem cert/ca-notifications-cert.pem ./services/notifications_service/configs/cert && \
	cp cert/ca-notifications-cert.pem cert/client-notifications-cert.pem cert/client-notifications-key.pem ./services/auth_service/configs/cert

gen_ssl_chat:
	cd cert &&  ./gen_service.sh chat; cd .. && \
	cp cert/server-chat-cert.pem cert/server-chat-key.pem cert/ca-chat-cert.pem ./services/chat_service/configs/cert
