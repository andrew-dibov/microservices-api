package clients

import (
	"microservices-api/internal/configs"
	"microservices-api/pkg/api/currency"

	"google.golang.org/grpc"
)

type CurrencyClient struct {
	grpc currency.CurrencyClient
	conn *grpc.ClientConn
	conf *configs.AppConfig
}
