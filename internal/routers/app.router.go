package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/middlewares"

	"log/slog"
	"net/http"
)

func NewAppRouter(hs *Handlers, cfg *configs.AppConfig, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/rate", hs.Curr.Rate)
	mux.HandleFunc("GET /api/v1/rates", hs.Curr.Rates)

	mux.HandleFunc("POST /api/v1/convert", hs.Conv.Convert)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	rtr := middlewares.Log(mux, log)
	rtr = middlewares.Recover(rtr, log)
	rtr = middlewares.Auth(rtr, log, cfg.Keys, cfg.Open)

	return rtr
}
