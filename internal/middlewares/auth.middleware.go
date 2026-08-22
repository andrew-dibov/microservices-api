package middlewares

import (
	"encoding/json"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"net/http"
)

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (middleware *AuthMiddleware) Middleware(next http.Handler, appConfig *configs.AppConfig, appLogger *loggers.AppLogger) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if appConfig.Security.OpenEndpoints[req.URL.Path] {
			next.ServeHTTP(res, req)
			return
		}

		apiKey := req.Header.Get("X-API-Key")

		if apiKey == "" {
			appLogger.Warn("AuthMiddleware : absent key",
				"path", req.URL.Path,
				"method", req.Method,
			)

			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)

			if err := json.NewEncoder(res).Encode(map[string]string{"error": "Absent key"}); err != nil {
				appLogger.Error("AuthMiddleware : json response failed", "error", err)
			}
			return
		}

		if !appConfig.Security.ApiKeys[apiKey] {
			appLogger.Warn("AuthMiddleware : wrong key",
				"api_key", apiKey,
				"path", req.URL.Path,
				"method", req.Method,
			)

			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)

			if err := json.NewEncoder(res).Encode(map[string]string{"error": "Wrong key"}); err != nil {
				appLogger.Error("AuthMiddleware : json response failed", "error", err)
			}
			return
		}

		next.ServeHTTP(res, req)
	})
}
