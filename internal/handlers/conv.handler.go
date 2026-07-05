package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"

	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

func NewConvHandler(c *clients.ConvClient, cfg *configs.AppConfig, log *slog.Logger) *ConvHandler {
	return &ConvHandler{
		c:   c,
		log: log,
		lim: ConvLimits{
			convert: rate.NewLimiter(rate.Limit(cfg.Limits.ConvertLimit), cfg.Limits.ConvertBurst),
		},
	}
}

func (h *ConvHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("json response failed",
			"error", err,
		)
	}
}

func (h *ConvHandler) Convert(w http.ResponseWriter, r *http.Request) {
	if !h.lim.convert.Allow() {
		h.log.Warn("limit exceeded")
		h.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var body ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "invalid body",
		})
		return
	}

	if body.FromCurrency == "" || body.ToCurrency == "" {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is empty",
		})
		return
	}

	if len(body.FromCurrency) != 3 || len(body.ToCurrency) != 3 {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "fromCurrency or toCurrency is not 3 chars",
		})
		return
	}

	for _, ch := range body.FromCurrency {
		if ch < 'A' || ch > 'Z' {
			h.res(w, http.StatusBadRequest, map[string]string{
				"error": "fromCurrency has invalid chars",
			})
			return
		}
	}

	for _, ch := range body.ToCurrency {
		if ch < 'A' || ch > 'Z' {
			h.res(w, http.StatusBadRequest, map[string]string{
				"error": "toCurrency has invalid chars",
			})
			return
		}
	}

	if body.Amount <= 0 {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "amount must be positive",
		})
		return
	}

	if body.Amount > 1e12 {
		h.res(w, http.StatusBadRequest, map[string]string{
			"error": "amount is too large",
		})
		return
	}

	ctx := r.Context()
	data, err := h.c.Convert(ctx, body.FromCurrency, body.ToCurrency, body.Amount)

	if err != nil {
		h.log.Error("conversion operation failed",
			"error", err,
		)
		h.res(w, http.StatusInternalServerError, map[string]string{
			"error": "conversion operation failed",
		})
		return
	}

	h.res(w, http.StatusOK, ConvertResponse{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Amount:       data.Amount,
		Result:       data.Result,
		Rate:         data.Rate,
	})
}
