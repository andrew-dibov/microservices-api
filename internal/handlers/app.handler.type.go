package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"microservices-api/internal/registries"
)

type AppHandler struct {
	appConfig          *configs.AppConfig
	appLogger          *loggers.AppLogger
	currencyClient     *clients.CurrencyClient
	conversionClient   *clients.ConversionClient
	prometheusRegistry *registries.PrometheusRegistry
}
