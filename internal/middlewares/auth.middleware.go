package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Auth(n http.Handler, l *slog.Logger, ks map[string]bool, op map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if op[r.URL.Path] {
			n.ServeHTTP(w, r)
			return
		}

		k := r.Header.Get("X-API-Key")

		if k == "" {
			l.Warn("absent key",
				"path", r.URL.Path,
				"method", r.Method,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Absent Key"}); err != nil {
				l.Error("json response failed",
					"error", err,
				)
			}
			return
		}

		if !ks[k] {
			l.Warn("wrong key",
				"path", r.URL.Path,
				"method", r.Method,
				"key", k,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Wrong Key"}); err != nil {
				l.Error("json response failed",
					"error", err,
				)
			}
			return
		}

		n.ServeHTTP(w, r)
	})
}
