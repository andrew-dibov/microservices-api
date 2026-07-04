package clients

import (
	"context"
	"fmt"
	"microservices-api/pkg/api/currency"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewCurrencyClient(addr string, tout time.Duration) (*CurrencyClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4*1024*1024), grpc.MaxCallSendMsgSize(4*1024*1024)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10 * time.Second, Timeout: 1 * time.Second}),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 2 * time.Second}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
		"methodConfig": [{ "name": [{"service": "currency.Currency"}],
			"retryPolicy": {
				"maxAttempts": 3,
				"maxBackoff": "1s",
				"backoffMultiplier": 2,
				"initialBackoff": "0.1s",
				"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`))

	if err != nil {
		return nil, fmt.Errorf("currency client : %w", err)
	}

	return &CurrencyClient{
		grpc: currency.NewCurrencyClient(conn),
		conn: conn,
		tout: tout,
	}, nil
}

func (cl *CurrencyClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}

func (cl *CurrencyClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (*currency.RateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, cl.tout)
	defer cancel()

	return cl.grpc.GetRate(ctx, &currency.RateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}

func (cl *CurrencyClient) GetAllRates(ctx context.Context, baseCurrency string) (*currency.RatesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, cl.tout)
	defer cancel()

	return cl.grpc.GetAllRates(ctx, &currency.RatesRequest{
		BaseCurrency: baseCurrency,
	})
}
