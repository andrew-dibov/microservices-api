package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func Log(nxt http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		came := time.Now()

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		nxt.ServeHTTP(sw, r)

		elapsed := time.Since(came).Milliseconds()

		log.Info("http request processed",
			"path", r.URL.Path, "method", r.Method, "status", sw.status, "elapsed", elapsed)
	})
}
