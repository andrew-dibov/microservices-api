package middlewares

import (
	"context"
	"microservices-api/internal/tools"
	"net/http"
)

func NewTraceMiddleware() *TraceMiddleware {
	return &TraceMiddleware{}
}

func (middleware *TraceMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := req.Header.Get("X-Request-ID")

		if id == "" {
			id = tools.GenID()
		}

		ctx := context.WithValue(req.Context(), tools.ID{}, id)
		req = req.WithContext(ctx)

		res.Header().Set("X-Request-ID", id)
		next.ServeHTTP(res, req)
	})
}
