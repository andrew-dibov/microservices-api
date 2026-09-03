package routers

import (
	"microservices-api/internal/handlers"
	"net/http"
)

type AppRouter struct {
	http.Handler
}

type AppRouterHandlers struct {
	App               *handlers.AppHandler
	CurrencyHandler   *handlers.CurrencyHandler
	ConversionHandler *handlers.ConversionHandler
}
