package handlers

import (
	"encoding/json"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"net/http"

	"golang.org/x/time/rate"
)

func NewConversionHandler(appConfig *configs.AppConfig, appLogger *loggers.AppLogger, conversionClient *clients.ConversionClient) *ConversionHandler {
	return &ConversionHandler{
		conversionClient: conversionClient,
		appLogger:        appLogger,
		conversionRateLimiters: ConversionRateLimiters{
			convert: rate.NewLimiter(rate.Limit(appConfig.ConversionService.Limits.Convert.Limit), appConfig.ConversionService.Limits.Convert.Burst),
		},
	}
}

/* --- --- --- */

func (handler *ConversionHandler) respond(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(data); err != nil {
		handler.appLogger.Error("ConversionHandler json response failed", "error", err)
	}
}

/* --- --- --- */

func (handler *ConversionHandler) Convert(res http.ResponseWriter, req *http.Request) {
	if !handler.conversionRateLimiters.convert.Allow() {
		handler.appLogger.Warn("Convert method limit exceeded")
		handler.respond(res, http.StatusTooManyRequests, map[string]string{"error": "Convert method limit exceeded"})
		return
	}

	/* --- --- --- */

	req.Body = http.MaxBytesReader(res, req.Body, 1024*1024)

	var convertRequest ConvertRequest
	if err := json.NewDecoder(req.Body).Decode(&convertRequest); err != nil {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	/* --- --- --- */

	if convertRequest.FromCurrency == "" || convertRequest.ToCurrency == "" {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency or toCurrency is empty"})
		return
	}

	if len(convertRequest.FromCurrency) != 3 || len(convertRequest.ToCurrency) != 3 {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency or toCurrency is not 3 chars"})
		return
	}

	for _, char := range convertRequest.FromCurrency {
		if char < 'A' || char > 'Z' {
			handler.respond(res, http.StatusBadRequest, map[string]string{"error": "fromCurrency has invalid chars"})
			return
		}
	}

	for _, char := range convertRequest.ToCurrency {
		if char < 'A' || char > 'Z' {
			handler.respond(res, http.StatusBadRequest, map[string]string{"error": "toCurrency has invalid chars"})
			return
		}
	}

	if convertRequest.Amount <= 0 {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "amount is not positive"})
		return
	}

	if convertRequest.Amount > 1e12 {
		handler.respond(res, http.StatusBadRequest, map[string]string{"error": "amount is too large"})
		return
	}

	/* --- --- --- */

	ctx := req.Context()
	data, err := handler.conversionClient.Convert(ctx, convertRequest.FromCurrency, convertRequest.ToCurrency, convertRequest.Amount)

	if err != nil {
		handler.appLogger.Error("ConversionHandler convert operation failed", "error", err)
		handler.respond(res, http.StatusInternalServerError, map[string]string{"error": "operation failed"})
		return
	}

	handler.respond(res, http.StatusOK, ConvertResponse{
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Amount:       data.Amount,
		Result:       data.Result,
		Rate:         data.Rate,
	})
}
