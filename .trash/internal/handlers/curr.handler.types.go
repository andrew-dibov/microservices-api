package handlers

import (
	"microservices-api/internal/clients"

	"log/slog"

	"golang.org/x/time/rate"
)

type CurrHandler struct {
	cl  *clients.CurrClient
	l   *slog.Logger
	lim CurrLim
}

type CurrLim struct {
	rate  *rate.Limiter
	rates *rate.Limiter
}

type RateRes struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
}

type RatesRes struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
}
