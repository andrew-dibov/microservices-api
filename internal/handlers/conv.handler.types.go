package handlers

import (
	"microservices-api/internal/clients"

	"log/slog"

	"golang.org/x/time/rate"
)

type ConvHandler struct {
	c   *clients.ConvClient
	log *slog.Logger
	lim ConvLimits
}

type ConvLimits struct {
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
