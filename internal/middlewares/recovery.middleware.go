package middlewares

import (
	"log/slog"
	"net/http"
)

func Recovery(nxt http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rc := recover(); rc != nil {
				log.Error("http server panicked",
					"path", r.URL.Path, "method", r.Method, "recover", rc)

				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error\n"))
			}
		}()

		nxt.ServeHTTP(w, r)
	})
}
