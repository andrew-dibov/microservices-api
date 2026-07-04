package middlewares

import (
	"log/slog"
	"net/http"
)

func Auth(nxt http.Handler, log *slog.Logger, kys map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			nxt.ServeHTTP(w, r)
			return
		}

		k := r.Header.Get("X-API-Key")

		if k == "" {
			log.Warn("absent api key", "path", r.URL.Path, "method", r.Method)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Absent API Key\n"))

			return
		}

		if !kys[k] {
			log.Warn("wrong api key", "path", r.URL.Path, "method", r.Method, "key", k)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Wrong API Key\n"))

			return
		}

		nxt.ServeHTTP(w, r)
	})
}
