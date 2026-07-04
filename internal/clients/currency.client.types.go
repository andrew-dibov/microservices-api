package clients

import (
	"microservices-api/pkg/api/currency"
	"time"

	"google.golang.org/grpc"
)

type CurrencyClient struct {
	grpc currency.CurrencyClient
	conn *grpc.ClientConn
	tout time.Duration
}
