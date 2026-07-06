package clients

import (
	"microservices-api/internal/middlewares"
	"microservices-api/pkg/api/currency"

	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

func NewCurrClient(addr string, tout time.Duration) (*CurrClient, error) {
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
		return nil, fmt.Errorf("failed to init : %w", err)
	}

	return &CurrClient{
		grpc: currency.NewCurrencyClient(conn),
		conn: conn,
		tout: tout,
	}, nil
}

func (cl *CurrClient) Health(ctx context.Context) error {
	state := cl.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return fmt.Errorf("connection state : %s", state)
	}
	return nil
}

func (cl *CurrClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}

func (cl *CurrClient) Rate(ctx context.Context, fromCurrency string, toCurrency string) (*currency.RateResponse, error) {
	if reqID := middlewares.GetReqID(ctx); reqID != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Request-ID", reqID)
	}

	ctx, can := context.WithTimeout(ctx, cl.tout)
	defer can()

	return cl.grpc.Rate(ctx, &currency.RateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}

func (cl *CurrClient) Rates(ctx context.Context, baseCurrency string) (*currency.RatesResponse, error) {
	if reqID := middlewares.GetReqID(ctx); reqID != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Request-ID", reqID)
	}

	ctx, can := context.WithTimeout(ctx, cl.tout)
	defer can()

	return cl.grpc.Rates(ctx, &currency.RatesRequest{
		BaseCurrency: baseCurrency,
	})
}
