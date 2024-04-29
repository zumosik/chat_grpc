APP_NAME=app

.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i public/css/input.css -o ./public/styles.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i public/css/input.css -o ./public/styles.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: dev
dev:
	go build -o ./tmp/$(APP_NAME) ./cmd/$(APP_NAME)/main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

.PHONY: test
test:
	  go test -race -v -timeout 30s ./...

.PHONY: install-tailwind-linux
install-tailwind-linux:
	curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
	chmod +x tailwindcss-linux-x64
	mv tailwindcss-linux-x64 tailwindcss
