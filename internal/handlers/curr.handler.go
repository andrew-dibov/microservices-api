package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"

	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

func NewCurrHandler(c *clients.CurrClient, cfg *configs.AppConfig, log *slog.Logger) *CurrHandler {
	return &CurrHandler{
		c:   c,
		log: log,
		lim: CurrLimits{
			rate:  rate.NewLimiter(rate.Limit(cfg.Limits.RateLimit), cfg.Limits.RateBurst),
			rates: rate.NewLimiter(rate.Limit(cfg.Limits.RatesLimit), cfg.Limits.RatesBurst),
		},
	}
}

func (h *CurrHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("json response failed",
			"error", err,
		)
	}
}

func (h *CurrHandler) Rate(w http.ResponseWriter, r *http.Request) {
	if !h.lim.rate.Allow() {
		h.log.Warn("limit exceeded")
		h.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
		})
		return
	}

	fromCurrency := r.URL.Query().Get("fromCurrency")
	toCurrency := r.URL.Query().Get("toCurrency")

	if fromCurrency == "" || toCurrency == "" {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is empty",
		})
		return
	}

	if len(fromCurrency) != 3 || len(toCurrency) != 3 {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range fromCurrency {
		if ch < 'A' || ch > 'Z' {
			h.res(w, http.StatusBadRequest, map[string]string{
				"error": "fromCurrency has invalid chars",
			})
			return
		}
	}

	for _, ch := range toCurrency {
		if ch < 'A' || ch > 'Z' {
			h.res(w, http.StatusBadRequest, map[string]string{
				"error": "toCurrency has invalid chars",
			})
			return
		}
	}

	ctx := r.Context()
	data, err := h.c.Rate(ctx, fromCurrency, toCurrency)

	if err != nil {
		h.log.Error("rate operation failed",
			"error", err,
		)
		h.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rate operation failed",
		})
		return
	}

	h.res(w, http.StatusOK, RateResponse{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Rate:         data.Rate,
	})
}

func (h *CurrHandler) Rates(w http.ResponseWriter, r *http.Request) {
	if !h.lim.rates.Allow() {
		h.log.Warn("limit exceeded")
		h.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
		})
		return
	}

	baseCurrency := r.URL.Query().Get("baseCurrency")

	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	if len(baseCurrency) != 3 {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "baseCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range baseCurrency {
		if ch < 'A' || ch > 'Z' {
			h.res(w, http.StatusBadRequest, map[string]string{
				"error": "baseCurrency has invalid chars",
			})
			return
		}
	}

	ctx := r.Context()
	data, err := h.c.Rates(ctx, baseCurrency)

	if err != nil {
		h.log.Error("rates operation failed",
			"error", err,
		)
		h.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rates operation failed",
		})
		return
	}

	h.res(w, http.StatusOK, RatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
