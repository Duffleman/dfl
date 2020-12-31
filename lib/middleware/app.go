package middleware

import (
	"context"
	"net/http"

	"dfl/lib/rpc"
)

var AppContext rpc.ContextKey = "app"

func MountApp(a interface{}, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), AppContext, a))
}

func MountAppMiddleware(a interface{}) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, MountApp(a, r))
		}

		return http.HandlerFunc(fn)
	}

}
