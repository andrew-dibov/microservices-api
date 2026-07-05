package clients

import (
	"microservices-api/pkg/api/conversion"

	"time"

	"google.golang.org/grpc"
)

type ConvClient struct {
	grpc conversion.ConversionClient
	conn *grpc.ClientConn
	tout time.Duration
}
