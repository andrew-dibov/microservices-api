package middlewares

import (
	"encoding/json"
	"microservices-api/internal/loggers"
	"net/http"
)

func NewRecoverMiddleware() *RecoverMiddleware {
	return &RecoverMiddleware{}
}

func (middleware *RecoverMiddleware) Middleware(next http.Handler, appLogger *loggers.AppLogger) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if recover := recover(); recover != nil {
				appLogger.Error("RecoverMiddleware : internal error",
					"recover", recover,
					"path", req.URL.Path,
					"method", req.Method,
				)

				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusInternalServerError)

				if err := json.NewEncoder(res).Encode(map[string]string{"error": "Internal server error"}); err != nil {
					appLogger.Error("RecoverMiddleware : json response failed", "error", err)
				}
			}
		}()

		next.ServeHTTP(res, req)
	})
}
