package handlers

import (
	"context"
	"encoding/json"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"microservices-api/internal/registries"
	"net/http"
)

func NewAppHandler(appConfig *configs.AppConfig, appLogger *loggers.AppLogger, currencyClient *clients.CurrencyClient, conversionClient *clients.ConversionClient, prometheusRegistry *registries.PrometheusRegistry) *AppHandler {
	return &AppHandler{
		appConfig: appConfig,
		appLogger: appLogger,

		currencyClient:     currencyClient,
		conversionClient:   conversionClient,
		prometheusRegistry: prometheusRegistry,
	}
}

/* --- --- --- */

func (handler *AppHandler) respond(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		handler.appLogger.Error("ConversionHandler json response failed", "error", err)
	}
}

/* --- --- --- */

func (handler *AppHandler) Root(res http.ResponseWriter, req *http.Request) {
	handler.respond(res, http.StatusOK, map[string]interface{}{"app": handler.appConfig.App.Name})
}

func (handler *AppHandler) Livez(res http.ResponseWriter, req *http.Request) {
	handler.respond(res, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (handler *AppHandler) Readyz(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), handler.appConfig.App.ReadyzTimeout)
	defer cancel()

	/* --- --- --- */

	services := make(map[string]string)
	status := "ok"

	if err := handler.currencyClient.Health(ctx); err != nil {
		services["currency"] = "unavailable"
		status = "degraded"
	} else {
		services["currency"] = "ok"
	}

	if err := handler.conversionClient.Health(ctx); err != nil {
		services["conversion"] = "unavailable"
		status = "degraded"
	} else {
		services["conversion"] = "ok"
	}

	result := map[string]interface{}{
		"services": services,
		"status":   status,
	}

	/* --- --- --- */

	if status != "ok" {
		handler.respond(res, http.StatusServiceUnavailable, result)
	} else {
		handler.respond(res, http.StatusOK, result)
	}
}

func (handler *AppHandler) Healthz(res http.ResponseWriter, req *http.Request) {
	handler.respond(res, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (handler *AppHandler) Metrics(res http.ResponseWriter, req *http.Request) {
	handler.prometheusRegistry.Metrics().ServeHTTP(res, req)
}
