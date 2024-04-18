APP_NAME = app

all: build run

.PHONY: build
build:
	go build -o $(APP_NAME) ./cmd/server/main.go

.PHONY: run
run: build
	./$(APP_NAME) --config=configs/local.yml

.PHONY: clean
clean:
	rm -f $(APP_NAME)

