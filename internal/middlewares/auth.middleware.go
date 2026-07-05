package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Auth(n http.Handler, log *slog.Logger, kys map[string]bool, opn map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opn[r.URL.Path] {
			n.ServeHTTP(w, r)
			return
		}

		k := r.Header.Get("X-API-Key")

		if k == "" {
			log.Warn("absent key",
				"path", r.URL.Path,
				"method", r.Method,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Absent Key"}); err != nil {
				log.Error("json response failed",
					"error", err,
				)
			}
			return
		}

		if !kys[k] {
			log.Warn("wrong key",
				"path", r.URL.Path,
				"method", r.Method,
				"key", k,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Wrong Key"}); err != nil {
				log.Error("json response failed",
					"error", err,
				)
			}
			return
		}

		n.ServeHTTP(w, r)
	})
}
