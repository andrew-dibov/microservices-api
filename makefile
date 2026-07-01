generate:
		protoc --go_out=. --go-grpc_out=. proto/history/history.proto
		protoc --go_out=. --go-grpc_out=. proto/currency/currency.proto
		protoc --go_out=. --go-grpc_out=. proto/conversion/conversion.proto

build:
		go build -o bin/microservices-api cmd/microservices-api/main.go

run:
		go run cmd/microservices-api/main.go
