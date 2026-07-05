package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Recover(n http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("internal error",
					"path", r.URL.Path,
					"method", r.Method,
					"recover", rec,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"}); err != nil {
					log.Error("json response failed",
						"error", err,
					)
				}
			}
		}()

		n.ServeHTTP(w, r)
	})
}
