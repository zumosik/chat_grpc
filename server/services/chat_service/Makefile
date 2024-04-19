APP_NAME = app

all: build_cl run_cl

.PHONY: build
build_cl:
	go build -o $(APP_NAME) ./cmd/test_tls_client/main.go

.PHONY: run
run_cl: build_cl
	./$(APP_NAME) --config=configs/local.yml

.PHONY: clean
clean_cl:
	rm -f $(APP_NAME)

