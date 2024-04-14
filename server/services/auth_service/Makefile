APP_NAME = app
POSTGRES_MIGARTE_DB=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

all: build run

.PHONY: build
build:
	go build -o $(APP_NAME) ./cmd/main.go

.PHONY: run
run: build
	./$(APP_NAME) --config=configs/local.yml

.PHONY: clean
clean:
	rm -f $(APP_NAME)

