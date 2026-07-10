package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"microservices-api/internal/clients"
	"microservices-api/internal/registries"
	"net/http"
	"time"
)

func NewAppHandler(curr *clients.CurrClient, conv *clients.ConvClient, preg *registries.PromRegistry, log *slog.Logger) *AppHandler {
	return &AppHandler{
		curr: curr,
		conv: conv,
		preg: preg,
		log:  log,
	}
}

func (h *AppHandler) Livez(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"}); err != nil {
		h.log.Error("json response failed",
			"error", err,
		)
	}
}

func (h *AppHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, can := context.WithTimeout(r.Context(), 2*time.Second)
	defer can()

	svcs := make(map[string]string)
	stat := "ok"

	if err := h.curr.Health(ctx); err != nil {
		svcs["currency"] = "unavailable"
		stat = "degraded"
	} else {
		svcs["currency"] = "ok"
	}

	if err := h.conv.Health(ctx); err != nil {
		svcs["conversion"] = "unavailable"
		stat = "degraded"
	} else {
		svcs["conversion"] = "ok"
	}

	res := map[string]interface{}{
		"status":   stat,
		"services": svcs,
	}

	w.Header().Set("Content-Type", "application/json")

	if stat != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.log.Error("json response failed",
			"error", err,
		)
	}
}

func (h *AppHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"}); err != nil {
		h.log.Error("json response failed",
			"error", err,
		)
	}
}

func (h *AppHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	h.preg.Metrics().ServeHTTP(w, r)
}
