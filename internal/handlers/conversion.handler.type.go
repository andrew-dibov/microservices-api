package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/loggers"

	"golang.org/x/time/rate"
)

type ConversionHandler struct {
	conversionClient       *clients.ConversionClient
	appLogger              *loggers.AppLogger
	conversionRateLimiters ConversionRateLimiters
}

type ConversionRateLimiters struct {
	convert *rate.Limiter
}

/* --- --- --- */

type ConvertRequest struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       float64 `json:"amount"`
}

type ConvertResponse struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       float64 `json:"amount"`
	Result       float64 `json:"result"`
	Rate         float64 `json:"rate"`
}
