package handlers

import (
	"microservices-api/internal/clients"

	"log/slog"

	"golang.org/x/time/rate"
)

type ConvHandler struct {
	cl  *clients.ConvClient
	l   *slog.Logger
	lim ConvLim
}

type ConvLim struct {
	convert *rate.Limiter
}

type ConvertReq struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       float64 `json:"amount"`
}

type ConvertRes struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       float64 `json:"amount"`
	Result       float64 `json:"result"`
	Rate         float64 `json:"rate"`
}
