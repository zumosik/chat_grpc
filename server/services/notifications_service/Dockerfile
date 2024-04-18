FROM golang:1.22.2-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/server/main.go


FROM alpine as runner

COPY --from=builder /usr/local/src/bin/app /app

ARG config_path_from=configs/prod.yaml
ENV CONFIG_PATH=config.yaml


COPY $config_path_from $CONFIG_PATH

CMD ["/app"]
