package middleware

import (
	"context"
	"net/http"
)

type CtxRequest struct{}

func HumaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), CtxRequest{}, r))
		next.ServeHTTP(w, r)
	})
}
