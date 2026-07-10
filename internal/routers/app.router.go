package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/middlewares"

	"log/slog"
	"net/http"
)

func NewAppRouter(hs *Handlers, cfg *configs.AppConfig, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /livez", hs.App.Livez)
	mux.HandleFunc("GET /readyz", hs.App.Readyz)
	mux.HandleFunc("GET /healthz", hs.App.Healthz)
	mux.HandleFunc("GET /metrics", hs.App.Metrics)

	mux.HandleFunc("GET /api/v1/rate", hs.Curr.Rate)
	mux.HandleFunc("GET /api/v1/rates", hs.Curr.Rates)
	mux.HandleFunc("POST /api/v1/convert", hs.Conv.Convert)

	rtr := middlewares.Auth(mux, log, cfg.Keys, cfg.Open)

	rtr = middlewares.Recover(rtr, log)
	rtr = middlewares.Log(rtr, log)
	rtr = middlewares.Trace(rtr)

	return rtr
}
