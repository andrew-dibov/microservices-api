package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/loggers"

	"golang.org/x/time/rate"
)

type CurrencyHandler struct {
	currencyClient       *clients.CurrencyClient
	appLogger            *loggers.AppLogger
	currencyRateLimiters CurrencyRateLimiters
}

type CurrencyRateLimiters struct {
	rate  *rate.Limiter
	rates *rate.Limiter
}

/* --- --- --- */

type RateResponse struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
}

type RatesResponse struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
}
