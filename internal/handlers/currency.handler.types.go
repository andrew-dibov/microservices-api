package handlers

import (
	"log/slog"
	"microservices-api/internal/clients"

	"golang.org/x/time/rate"
)

type CurrencyHandler struct {
	cl  *clients.CurrencyClient
	log *slog.Logger
	lim CurrencyLimits
}

type CurrencyLimits struct {
	rate  *rate.Limiter
	rates *rate.Limiter
}

type RateResponse struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
}

type RatesResponse struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
}
