package handlers

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"

	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

func NewConvHandler(cl *clients.ConvClient, c *configs.AppConfig, l *slog.Logger) *ConvHandler {
	return &ConvHandler{
		cl: cl,
		l:  l,
		lim: ConvLim{
			convert: rate.NewLimiter(rate.Limit(c.Limt.ConvertLim), c.Limt.ConvertBur),
		},
	}
}

func (ha *ConvHandler) res(w http.ResponseWriter, stat int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ha.l.Error("json response failed",
			"error", err,
		)
	}
}

func (ha *ConvHandler) Convert(w http.ResponseWriter, r *http.Request) {
	if !ha.lim.convert.Allow() {
		ha.l.Warn("limit exceeded")
		ha.res(w, http.StatusTooManyRequests, map[string]string{
			"message": "limit exceeded",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var body ConvertReq
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
		ha.l.Error("conversion operation failed",
			"error", err,
		)
		ha.res(w, http.StatusInternalServerError, map[string]string{
			"error": "conversion operation failed",
		})
		return
	}

	ha.res(w, http.StatusOK, ConvertRes{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Amount:       data.Amount,
		Result:       data.Result,
		Rate:         data.Rate,
	})
}
