package middlewares

import (
	"context"
	"microservices-api/internal/modules"
	"net/http"
)

func NewTraceMiddleware() *TraceMiddleware {
	return &TraceMiddleware{}
}

func (middleware *TraceMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := req.Header.Get("x-request-id")

		if id == "" {
			id = modules.GenID()
		}

		ctx := context.WithValue(req.Context(), modules.ID{}, id)
		req = req.WithContext(ctx)

		res.Header().Set("x-request-id", id)
		next.ServeHTTP(res, req)
	})
}
