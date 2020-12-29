package middleware

import (
	"context"
	"net/http"

	"dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
)

var authenticatedPaths = map[string]struct{}{
	"/add_shortcut":       {},
	"/created_signed_url": {},
	"/delete_resource":    {},
	"/list_resources":     {},
	"/remove_shortcut":    {},
	"/resave_hashes":      {},
	"/set_nsfw":           {},
	"/shorten_url":        {},
	"/upload_file":        {},
	"/view_details":       {},
}

func AuthMiddleware(users map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := authenticatedPaths[r.URL.Path]; !ok {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
				return
			}

			for username, key := range users {
				if key == authHeader {
					ctx := context.WithValue(r.Context(), auth.UserContextKey, username)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
			return
		}

		return http.HandlerFunc(fn)
	}
}
