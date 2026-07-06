package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func Log(n http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		came := time.Now()

		sw := &statWriter{ResponseWriter: w, status: http.StatusOK}
		reqID := GetReqID(r.Context())

		n.ServeHTTP(sw, r)

		processed := time.Since(came).Milliseconds()

		log.Info("http request",
			"id", reqID,
			"path", r.URL.Path,
			"method", r.Method,
			"status", sw.status,
			"processed", processed,
		)
	})
}
