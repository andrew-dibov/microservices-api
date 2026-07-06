package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func Trace(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = genReqID()
		}

		ctx := context.WithValue(r.Context(), ctxKeyReqID{}, id)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", id)
		n.ServeHTTP(w, r)
	})
}

func genReqID() string {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "0000000000000000"
	}

	return hex.EncodeToString(b)
}

func GetReqID(ctx context.Context) string {
	if v := ctx.Value(ctxKeyReqID{}); v != nil {
		return v.(string)
	}
	return ""
}
