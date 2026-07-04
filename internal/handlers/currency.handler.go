package handlers

import (
	"encoding/json"
	"log/slog"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"net/http"

	"golang.org/x/time/rate"
)

func NewCurrencyHandler(cl *clients.CurrencyClient, cfg *configs.AppConfig, log *slog.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		cl:  cl,
		log: log,
		lim: CurrencyLimits{
			rate:  rate.NewLimiter(rate.Limit(cfg.Limits.RateLimit), cfg.Limits.RateBurst),
			rates: rate.NewLimiter(rate.Limit(cfg.Limits.RatesLimit), cfg.Limits.RatesBurst),
		},
	}
}

func (ha *CurrencyHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ha.log.Error("json response failed", "error", err)
	}
}

func (ha *CurrencyHandler) Rate(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.rate.Allow() {
		ha.log.Warn("rate limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "rate limit exceeded",
		})
		return
	}

	fromCurrency := r.URL.Query().Get("fromCurrency")
	toCurrency := r.URL.Query().Get("toCurrency")

	if fromCurrency == "" || toCurrency == "" {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is empty",
		})
		return
	}

	if len(fromCurrency) != 3 || len(toCurrency) != 3 {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range fromCurrency {
		if ch < 'A' || ch > 'Z' {
			ha.res(w, http.StatusBadRequest, map[string]string{
				"error": "fromCurrency has invalid chars",
			})
			return
		}
	}

	for _, ch := range toCurrency {
		if ch < 'A' || ch > 'Z' {
			ha.res(w, http.StatusBadRequest, map[string]string{
				"error": "toCurrency has invalid chars",
			})
			return
		}
	}

	ctx := r.Context()
	data, err := ha.cl.GetRate(ctx, fromCurrency, toCurrency)

	if err != nil {
		ha.log.Error("rate operation failed", "error", err)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rate operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, RateResponse{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Rate:         data.Rate,
	})
}

func (ha *CurrencyHandler) Rates(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.rates.Allow() {
		ha.log.Warn("rate limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "rate limit exceeded",
		})
		return
	}

	baseCurrency := r.URL.Query().Get("baseCurrency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	if len(baseCurrency) != 3 {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "baseCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range baseCurrency {
		if ch < 'A' || ch > 'Z' {
			ha.res(w, http.StatusBadRequest, map[string]string{
				"error": "baseCurrency has invalid chars",
			})
			return
		}
	}

	ctx := r.Context()
	data, err := ha.cl.GetAllRates(ctx, baseCurrency)

	if err != nil {
		ha.log.Error("rates operation failed", "error", err)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rates operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, RatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
