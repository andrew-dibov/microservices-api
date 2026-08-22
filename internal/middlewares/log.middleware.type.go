package middlewares

import "net/http"

type LogMiddleware struct{}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (writer *responseWriter) WriteHeader(code int) {
	writer.ResponseWriter.WriteHeader(code)
	writer.status = code
}
