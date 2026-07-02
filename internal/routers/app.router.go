package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/handlers"
	"microservices-api/internal/middlewares"

	"log/slog"
	"net/http"
)

func NewAppRouter(cfg *configs.AppConfig, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/rate", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /api/v1/rates", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /api/v1/convert", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("GET /api/v1/health", handlers.Health)

	rtr := middlewares.Log(mux, log)
	rtr = middlewares.Recovery(rtr, log)

	return rtr
}
