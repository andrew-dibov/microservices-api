package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/middlewares"

	"log/slog"
	"net/http"
)

func NewAppRouter(hds *Handlers, cfg *configs.AppConfig, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/rate", hds.Currency.Rate)
	mux.HandleFunc("GET /api/v1/rates", hds.Currency.Rates)
	mux.HandleFunc("POST /api/v1/convert", hds.Conversion.Convert)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	rtr := middlewares.Log(mux, log)
	rtr = middlewares.Recovery(rtr, log)

	return rtr
}
