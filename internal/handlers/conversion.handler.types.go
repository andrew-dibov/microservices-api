package handlers

import (
	"log/slog"
	"microservices-api/internal/clients"

	"golang.org/x/time/rate"
)

type ConversionHandler struct {
	cl  *clients.ConversionClient
	log *slog.Logger
	lim ConversionLimits
}

type ConversionLimits struct {
	convert *rate.Limiter
}

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
