package middlewares

import (
	"microservices-api/internal/loggers"
	"microservices-api/internal/tools"
	"net/http"
	"time"
)

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (middleware *LogMiddleware) Middleware(next http.Handler, appLogger *loggers.AppLogger) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqNow := time.Now()
		newRes := &responseWriter{ResponseWriter: res, status: http.StatusOK}

		reqID := tools.GetID(req.Context())
		next.ServeHTTP(newRes, req)

		reqProcessed := time.Since(reqNow).Milliseconds()

		appLogger.Info("LogMiddleware : processed",
			"id", reqID,
			"path", req.URL.Path,
			"method", req.Method,
			"status", newRes.status,
			"processed", reqProcessed,
		)
	})
}
