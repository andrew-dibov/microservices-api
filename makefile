.PHONY: install run start stop up down test clean help
.DEFAULT_GOAL := help

-include .env
export

# ---

install:
	@go mod tidy

	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
	@protoc --go_out=. --go-grpc_out=. proto/currency/currency.proto
	@protoc --go_out=. --go-grpc_out=. proto/conversion/conversion.proto

	@mkdir -p certs
	@[ -f certs/cert.pem ] || (cd certs && openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost")
	
	@docker image build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg CMD_NAME=${CMD_NAME} \
		--build-arg APP_PORT=${APP_PORT} \
		-t ${CMD_NAME}:latest .

run:
	@go run cmd/${CMD_NAME}/main.go

start:
	@docker container run -d --rm \
		--env-file .env \
		-e SECURITY_CERTIFICATE=/certs/cert.pem \
		-e SECURITY_KEY=/certs/key.pem \
		-p ${HOST_PORT}:${APP_PORT} \
		-v ./certs:/certs:ro \
		--name ${CMD_NAME} \
		${CMD_NAME}:latest
	@docker container ps --filter name=${CMD_NAME}

stop:
	@docker container stop ${CMD_NAME} > /dev/null 2>&1

up:
	@docker compose build
	@docker compose up -d --no-build

down:
	@docker compose down -v

test:
	@bash shell/test.sh

clean:
	@rm -rf certs

help:
	@echo "install : установить зависимости и собрать docker образ"
	@echo "run     : запустить приложение"
	@echo "start   : запустить docker контейнер"
	@echo "stop    : остановить docker контейнер"
	@echo "up      : запустить docker compose проект"
	@echo "down    : остановить docker compose проект"
	@echo "test    : запустить тестирование"
	@echo "clean   : удалить сертификаты"
