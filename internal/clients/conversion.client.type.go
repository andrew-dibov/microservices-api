package clients

import (
	"microservices-api/internal/configs"
	"microservices-api/pkg/api/conversion"

	"google.golang.org/grpc"
)

type ConversionClient struct {
	grpc conversion.ConversionClient
	conn *grpc.ClientConn
	conf *configs.AppConfig
}
