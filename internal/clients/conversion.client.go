package clients

import (
	"fmt"
	"microservices-api/pkg/api/conversion"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ConversionClient struct {
	cl   conversion.ConversionClient
	conn *grpc.ClientConn
	tout time.Duration
}

func NewConversionClient(addr string, tout time.Duration) (*ConversionClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4*1024*1024), grpc.MaxCallSendMsgSize(4*1024*1024)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10 * time.Second, Timeout: 1 * time.Second}),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 2 * time.Second}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
		"methodConfig": [{ "name": [{"service": "conversion.Conversion"}],
			"retryPolicy": {
				"maxAttempts": 3,
				"maxBackoff": "1s",
				"backoffMultiplier": 2,
				"initialBackoff": "0.1s",
				"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`),
	)

	if err != nil {
		return nil, fmt.Errorf("conversion client failed : %w", err)
	}

	// con.Connect()

	return &ConversionClient{
		cl:   conversion.NewConversionClient(conn),
		conn: conn,
		tout: tout,
	}, nil
}

func (cl *ConversionClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}
