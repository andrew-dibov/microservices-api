package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/middlewares"

	"log/slog"
	"net/http"
)

func NewAppRouter(hs *Handlers, c *configs.AppConfig, l *slog.Logger) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("GET /livez", hs.App.Livez)
	m.HandleFunc("GET /readyz", hs.App.Readyz)
	m.HandleFunc("GET /healthz", hs.App.Healthz)
	m.HandleFunc("GET /metrics", hs.App.Metrics)

	m.HandleFunc("GET /api/v1/rate", hs.Curr.Rate)
	m.HandleFunc("GET /api/v1/rates", hs.Curr.Rates)
	m.HandleFunc("POST /api/v1/convert", hs.Conv.Convert)

	r := middlewares.Auth(m, l, c.App.Keys, c.App.Open)

	r = middlewares.Recover(r, l)
	r = middlewares.Log(r, l)
	r = middlewares.Trace(r)

	return r
}
