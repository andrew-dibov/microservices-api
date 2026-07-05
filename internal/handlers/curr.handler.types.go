package handlers

import (
	"microservices-api/internal/clients"

	"log/slog"

	"golang.org/x/time/rate"
)

type CurrHandler struct {
	c   *clients.CurrClient
	log *slog.Logger
	lim CurrLimits
}

type CurrLimits struct {
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
