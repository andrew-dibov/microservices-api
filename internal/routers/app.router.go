package routers

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"microservices-api/internal/middlewares"
	"net/http"
)

func NewAppRouter(handlers *AppRouterHandlers, appConfig *configs.AppConfig, appLogger *loggers.AppLogger) *AppRouter {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /", )

	mux.HandleFunc("GET /livez", handlers.App.Livez)
	mux.HandleFunc("GET /readyz", handlers.App.Readyz)
	mux.HandleFunc("GET /healthz", handlers.App.Healthz)
	mux.HandleFunc("GET /metrics", handlers.App.Metrics)

	mux.HandleFunc("GET /api/v1/rate", handlers.CurrencyHandler.Rate)
	mux.HandleFunc("GET /api/v1/rates", handlers.CurrencyHandler.Rates)

	mux.HandleFunc("POST /api/v1/convert", handlers.ConversionHandler.Convert)

	router := middlewares.NewAuthMiddleware().Middleware(mux, appConfig, appLogger)
	router = middlewares.NewRecoverMiddleware().Middleware(router, appLogger)
	router = middlewares.NewLogMiddleware().Middleware(router, appLogger)
	router = middlewares.NewTraceMiddleware().Middleware(router)

	return &AppRouter{Handler: router}
}
