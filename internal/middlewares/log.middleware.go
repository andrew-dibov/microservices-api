package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

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
