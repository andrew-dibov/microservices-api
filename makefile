BIN_NAME := app
CMD_NAME := microservices-api

BIN_DIR := bin
TLS_DIR := certs
API_KEY := test-25

DOCKER_IMG := $(CMD_NAME):latest

.DEFAULT_GOAL := help
.PHONY: get_dependencies get_protoc generate_proto generate_certificates build_binary run_binary run_application build_docker_container run_docker_container stop_docker_container run_test remove_certificates remove_binary clean help

# ---

get_dependencies:
	@go mod tidy

get_protoc:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate_proto:
	@protoc --go_out=. --go-grpc_out=. proto/currency/currency.proto
	@protoc --go_out=. --go-grpc_out=. proto/conversion/conversion.proto

generate_certificates:
	@mkdir -p $(TLS_DIR)
	@[ -f $(TLS_DIR)/cert.pem ] || (cd $(TLS_DIR) && openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost")

# ---

build_binary: get_dependencies get_protoc generate_proto generate_certificates
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) cmd/$(CMD_NAME)/main.go

run_binary: get_dependencies get_protoc generate_proto generate_certificates build_binary
	@SECURITY_CERTIFICATE=./certs/cert.pem SECURITY_KEY=./certs/key.pem ./$(BIN_DIR)/$(BIN_NAME)

run_application:
	@SECURITY_CERTIFICATE=./certs/cert.pem SECURITY_KEY=./certs/key.pem go run cmd/$(CMD_NAME)/main.go

# ---

build_docker_container:
	@docker image build -t $(DOCKER_IMG) .
	@echo ""

run_docker_container: get_dependencies get_protoc generate_proto generate_certificates build_docker_container
	docker container run -d --rm --name $(CMD_NAME) \
		-v ./certs:/certs:ro \
		-e SECURITY_CERTIFICATE=/certs/cert.pem \
		-e SECURITY_KEY=/certs/key.pem \
		-e SECURITY_API_KEYS=$(API_KEY) \
		-e APP_PROD=true \
		-p 8080:8080 $(DOCKER_IMG)

stop_docker_container:
	@docker ps -q --filter "name=$(CMD_NAME)" | xargs docker stop

# ---

run_test: get_dependencies get_protoc generate_proto generate_certificates build_docker_container
	docker container run -d --rm --name $(CMD_NAME)-test \
		-v ./certs:/certs:ro \
		-e SECURITY_CERTIFICATE=/certs/cert.pem \
		-e SECURITY_KEY=/certs/key.pem \
		-e SECURITY_API_KEYS=$(API_KEY) \
		-e APP_PROD=true \
		-p 8080:8080 $(DOCKER_IMG)
	@echo ""

	@sleep 2
	@chmod +x tests/*
	API_KEY=$(API_KEY) ./tests/basic_test.sh
	
	@echo ""
	docker container stop $(CMD_NAME)-test > /dev/null 2>&1

# ---

remove_certificates:
	@rm -rf $(TLS_DIR)

remove_binary:
	@rm -rf $(BIN_DIR)

clean:
	@rm -rf $(BIN_DIR) $(TLS_DIR)

# ---

help: