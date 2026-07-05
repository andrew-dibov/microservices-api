package clients

import (
	"microservices-api/pkg/api/conversion"

	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewConvClient(addr string, tout time.Duration) (*ConvClient, error) {
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
		}`))

	if err != nil {
		return nil, fmt.Errorf("failed to init : %w", err)
	}

	return &ConvClient{
		grpc: conversion.NewConversionClient(conn),
		conn: conn,
		tout: tout,
	}, nil
}

func (cl *ConvClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}

func (cl *ConvClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*conversion.ConvertResponse, error) {
	ctx, can := context.WithTimeout(ctx, cl.tout)
	defer can()

	return cl.grpc.Convert(ctx, &conversion.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
