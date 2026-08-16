package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"

	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

func NewCurrHandler(cl *clients.CurrClient, c *configs.AppConfig, l *slog.Logger) *CurrHandler {
	return &CurrHandler{
		cl: cl,
		l:  l,
		lim: CurrLim{
			rate:  rate.NewLimiter(rate.Limit(c.Limt.RateLim), c.Limt.RateBur),
			rates: rate.NewLimiter(rate.Limit(c.Limt.RatesLim), c.Limt.RatesBur),
		},
	}
}

func (ha *CurrHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ha.l.Error("json response failed",
			"error", err,
		)
	}
}

func (ha *CurrHandler) Rate(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.rate.Allow() {
		ha.l.Warn("limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
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
	data, err := ha.cl.Rate(ctx, fromCurrency, toCurrency)

	if err != nil {
		ha.l.Error("rate operation failed",
			"error", err,
		)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rate operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, RateRes{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Rate:         data.Rate,
	})
}

func (ha *CurrHandler) Rates(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.rates.Allow() {
		ha.l.Warn("limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
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
	data, err := ha.cl.Rates(ctx, baseCurrency)

	if err != nil {
		ha.l.Error("rates operation failed",
			"error", err,
		)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "rates operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, RatesRes{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
