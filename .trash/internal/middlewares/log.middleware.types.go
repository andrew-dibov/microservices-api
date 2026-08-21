package middlewares

import "net/http"

type statWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}
