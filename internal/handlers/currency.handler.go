package handlers

import (
	"encoding/json"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"net/http"

	"golang.org/x/time/rate"
)

func NewCurrencyHandler(appConfig *configs.AppConfig, appLogger *loggers.AppLogger, currencyClient *clients.CurrencyClient) *CurrencyHandler {
	return &CurrencyHandler{
		currencyClient: currencyClient,
		appLogger:      appLogger,
		currencyRateLimiters: CurrencyRateLimiters{
			rate:  rate.NewLimiter(rate.Limit(appConfig.CurrencyService.Limits.Rate.Limit), appConfig.CurrencyService.Limits.Rate.Burst),
			rates: rate.NewLimiter(rate.Limit(appConfig.CurrencyService.Limits.Rates.Limit), appConfig.CurrencyService.Limits.Rates.Burst),
		},
	}
}

/* --- --- --- */

func (handler *CurrencyHandler) respond(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		handler.appLogger.Error("CurrencyHandler json response failed", "error", err)
	}
}

/* --- --- --- */

func (handler *CurrencyHandler) Rate(res http.ResponseWriter, req *http.Request) {
	if !handler.currencyRateLimiters.rate.Allow() {
		handler.appLogger.Warn("Rate method limit exceeded")
		handler.respond(res, http.StatusTooManyRequests, map[string]string{"error": "Rate method limit exceeded"})
		return
	}

	/* --- --- --- */

	fromCurrency := req.URL.Query().Get("fromCurrency")
	toCurrency := req.URL.Query().Get("toCurrency")

	if fromCurrency == "" || toCurrency == "" {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency or toCurrency is empty"})
		return
	}

	if len(fromCurrency) != 3 || len(toCurrency) != 3 {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency or toCurrency is not 3 chars"})
		return
	}

	for _, char := range fromCurrency {
		if char < 'A' || char > 'Z' {
			handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency has invalid chars"})
			return
		}
	}

	for _, char := range toCurrency {
		if char < 'A' || char > 'Z' {
			handler.respond(res, http.StatusBadRequest, map[string]string{"error": "toCurrency has invalid chars"})
			return
		}
	}

	/* --- --- --- */

	ctx := req.Context()
	data, err := handler.currencyClient.Rate(ctx, fromCurrency, toCurrency)

	if err != nil {
		handler.appLogger.Error("CurrencyHandler rate operation failed", "error", err)
		handler.respond(res, http.StatusInternalServerError, map[string]string{"error": "operation failed"})
		return
	}

	handler.respond(res, http.StatusOK, RateResponse{FromCurrency: data.FromCurrency, ToCurrency: data.ToCurrency, Rate: data.Rate})
}

func (handler *CurrencyHandler) Rates(res http.ResponseWriter, req *http.Request) {
	if !handler.currencyRateLimiters.rates.Allow() {
		handler.appLogger.Warn("Rates method limit exceeded")
		handler.respond(res, http.StatusTooManyRequests, map[string]string{"error": "Rates method limit exceeded"})
		return
	}

	/* --- --- --- */

	baseCurrency := req.URL.Query().Get("baseCurrency")

	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	if len(baseCurrency) != 3 {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "baseCurrency is not 3 chars"})
		return
	}

	for _, char := range baseCurrency {
		if char < 'A' || char > 'Z' {
			handler.respond(res, http.StatusBadRequest, map[string]string{"error": "baseCurrency has invalid chars"})
			return
		}
	}

	/* --- --- --- */

	ctx := req.Context()
	data, err := handler.currencyClient.Rates(ctx, baseCurrency)

	if err != nil {
		handler.appLogger.Error("CurrencyHandler rates operation failed", "error", err)
		handler.respond(res, http.StatusInternalServerError, map[string]string{"error": "operation failed"})
		return
	}

	handler.respond(res, http.StatusOK, RatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
