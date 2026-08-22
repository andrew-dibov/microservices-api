package clients

import (
	"fmt"
	"microservices-api/internal/configs"
	"microservices-api/pkg/api/conversion"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewConversionClient(appConfig *configs.AppConfig) (*ConversionClient, error) {
	connection, err := grpc.NewClient(appConfig.ConversionService.Address,
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
		return nil, fmt.Errorf("failed to create ConversionClient : %w", err)
	}

	return &ConversionClient{
		grpc: conversion.NewConversionClient(connection),
		conn: connection,
		conf: appConfig,
	}, nil
}
