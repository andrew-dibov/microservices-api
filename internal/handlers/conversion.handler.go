package handlers

import (
	"encoding/json"
	"log/slog"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"net/http"

	"golang.org/x/time/rate"
)

func NewConversionHandler(cl *clients.ConversionClient, cfg *configs.AppConfig, log *slog.Logger) *ConversionHandler {
	return &ConversionHandler{
		cl:  cl,
		log: log,
		lim: ConversionLimits{
			convert: rate.NewLimiter(rate.Limit(cfg.Limits.ConvertLimit), cfg.Limits.ConvertBurst),
		},
	}
}

func (ha *ConversionHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ha.log.Error("json response failed", "error", err)
	}
}

func (ha *ConversionHandler) Convert(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.convert.Allow() {
		ha.log.Warn("rate limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "rate limit exceeded",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var body ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "invalid body",
		})
		return
	}

	if body.FromCurrency == "" || body.ToCurrency == "" {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is empty",
		})
		return
	}

	if len(body.FromCurrency) != 3 || len(body.ToCurrency) != 3 {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range body.FromCurrency {
		if ch < 'A' || ch > 'Z' {
			ha.res(w, http.StatusBadRequest, map[string]string{
				"error": "fromCurrency has invalid chars",
			})
			return
		}
	}

	for _, ch := range body.ToCurrency {
		if ch < 'A' || ch > 'Z' {
			ha.res(w, http.StatusBadRequest, map[string]string{
				"error": "toCurrency has invalid chars",
			})
			return
		}
	}

	if body.Amount <= 0 {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "amount must be positive",
		})
		return
	}

	if body.Amount > 1e12 {
		ha.res(w, http.StatusBadRequest, map[string]string{
			"error": "amount is too large",
		})
		return
	}

	ctx := r.Context()
	data, err := ha.cl.Convert(ctx, body.FromCurrency, body.ToCurrency, body.Amount)

	if err != nil {
		ha.log.Error("conversion operation failed", "error", err)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "conversion operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, ConvertResponse{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Amount:       data.Amount,
		Result:       data.Result,
		Rate:         data.Rate,
	})
}
