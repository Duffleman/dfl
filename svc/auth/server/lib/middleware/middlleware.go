package middleware

import (
	"context"
	"net/http"
	"strings"

	"dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"

	"github.com/dgrijalva/jwt-go"
)

var unauthenticatedPaths = map[string]struct{}{
	"/authorize":       {},
	"/get_public_cert": {},
	"/login":           {},
	"/register":        {},
	"/token":           {},
}

type DFLClaims struct {
	Scope    string `json:"scope"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthMiddleware(publicKey interface{}) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := unauthenticatedPaths[r.URL.Path]; ok {
				h.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
				return
			}

			parts := strings.Fields(authHeader)

			if len(parts) != 2 {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
				return
			}

			if parts[0] != "Bearer" {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
				return
			}

			tokenString := parts[1]

			var dflclaims DFLClaims

			token, err := jwt.ParseWithClaims(tokenString, &dflclaims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, cher.New("unexpected_signing_method", cher.M{
						"alg": token.Header["alg"],
					})
				}

				return publicKey, nil
			})
			if err != nil {
				rpc.HandleError(w, r, err)
				return
			}

			if !token.Valid {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
				return
			}

			claims := token.Claims.(*DFLClaims)

			authUser := auth.AuthUser{
				UserID:   claims.Id,
				Username: claims.Username,
				Scopes:   claims.Scope,
			}

			if !authUser.Can("dflauth:login") {
				rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserContextKey, authUser)

			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
