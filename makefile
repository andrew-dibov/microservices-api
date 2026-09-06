CMD_NAME := microservices-api

BIN_NAME := $(CMD_NAME)-bin
IMG_NAME := $(CMD_NAME):latest

API_KEY := $(CMD_NAME)-key

.DEFAULT_GOAL := help
.PHONY: get_deps get_protoc gen_protos gen_certs build_bin run_bin run_app build_container run_container stop_container test rm_certs rm_bin clean help

# ---

get_deps:
	@ go mod tidy

get_protoc:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen_protos:
	@protoc --go_out=. --go-grpc_out=. proto/currency/currency.proto
	@protoc --go_out=. --go-grpc_out=. proto/conversion/conversion.proto

gen_certs:
	@mkdir -p certs
	@[ -f certs/cert.pem ] || (cd certs && openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost")

# ---

build_bin: get_deps get_protoc gen_protos gen_certs
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/$(BIN_NAME) cmd/$(CMD_NAME)/main.go

run_bin: build_bin
	@SECURITY_CERTIFICATE=./certs/cert.pem SECURITY_KEY=./certs/key.pem ./bin/$(BIN_NAME)

# ---

run_app: get_deps get_protoc gen_protos gen_certs
	@SECURITY_CERTIFICATE=./certs/cert.pem SECURITY_KEY=./certs/key.pem go run cmd/$(CMD_NAME)/main.go

# ---

build_container:
	@docker image build -t $(IMG_NAME) . > /dev/null 2>&1

run_container: build_container
	@docker container run -d --rm --name $(CMD_NAME) \
		-v ./certs:/certs:ro \
		-e SECURITY_CERTIFICATE=/certs/cert.pem \
		-e SECURITY_KEY=/certs/key.pem \
		-e SECURITY_API_KEYS=$(API_KEY) \
		-e APP_PROD=true \
		-p 8080:8080 $(IMG_NAME)

stop_container:
	@docker container stop $(CMD_NAME) > /dev/null 2>&1

# ---

test: get_deps get_protoc gen_protos gen_certs build_container
	@docker container run -d --rm --name $(CMD_NAME)-test \
		-v ./certs:/certs:ro \
		-e SECURITY_CERTIFICATE=/certs/cert.pem \
		-e SECURITY_KEY=/certs/key.pem \
		-e SECURITY_API_KEYS=$(API_KEY) \
		-e APP_PROD=true \
		-p 2525:8080 $(IMG_NAME) > /dev/null 2>&1
	
	@sleep 5
	@chmod +x ./test.sh

	@API_KEY=$(API_KEY) ./test.sh || (docker stop $(CMD_NAME)-test > /dev/null 2>&1 && exit 1)
	@docker container stop $(CMD_NAME)-test > /dev/null 2>&1

# ---

rm_certs:
	@rm -rf certs

rm_bin:
	@rm -rf bin

clean:
	@rm -rf bin certs

# ---

help:
	@echo "get_deps        установить зависимости"
	@echo "get_protoc      установить protoc"
	@echo "gen_protos      сгенерировать .proto"
	@echo "gen_certs       сгенерировать сертификаты"
	@echo "build_bin       собрать бинарник"
	@echo "run_bin         запустить бинарник"
	@echo "run_app         запустить приложение"
	@echo "build_container собрать образ"
	@echo "run_container   запустить контейнер"
	@echo "stop_container  остановить контейнер"
	@echo "test            провести тесты"
	@echo "rm_certs        удалить сертификаты"
	@echo "rm_bin          удалить бинарник"
	@echo "clean           удалить сертификаты и бинарник"
	@echo "help            показать справку"