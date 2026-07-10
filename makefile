BIN_NAME := app
CMD_NAME := microservices-api

BIN_DIR := bin
TLS_DIR := certs

.DEFAULT_GOAL := help

generate:
	protoc --go_out=. --go-grpc_out=. proto/currency/currency.proto
	protoc --go_out=. --go-grpc_out=. proto/conversion/conversion.proto

	@mkdir -p $(TLS_DIR)
	@[ -f $(TLS_DIR)/cert.pem ] || (cd $(TLS_DIR) && openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost")

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) cmd/$(CMD_NAME)/main.go

run:
	go run cmd/$(CMD_NAME)/main.go

clean:
	rm -rf $(BIN_DIR) $(TLS_DIR)

help:
	@echo "make generate : сгенерировать protobuf и сертификаты"
	@echo "make build    : собрать бинарник"
	@echo "make run      : запустить сервер"
	@echo "make clean    : удалить bin и certs"

.PHONY: generate build run clean help