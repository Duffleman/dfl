package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	dfljwt "dfl/lib/jwt"
	"dfl/lib/rpc"

	"github.com/dgrijalva/jwt-go"
)

type HTTPResource struct {
	Verb  *string
	Path  *string
	Regex *regexp.Regexp
}

func AuthMiddleware(publicKey interface{}, bypassPaths []HTTPResource) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, path := range bypassPaths {
				if path.Verb != nil && *path.Verb != r.Method {
					continue
				}

				if path.Path != nil && *path.Path != r.URL.Path {
					continue
				}

				if path.Regex != nil && !path.Regex.MatchString(r.URL.Path) {
					continue
				}

				h.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil), nil)
				return
			}

			parts := strings.Fields(authHeader)

			if len(parts) != 2 {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil), nil)
				return
			}

			if parts[0] != "Bearer" {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil), nil)
				return
			}

			tokenString := parts[1]

			var dflclaims dfljwt.DFLClaims

			token, err := jwt.ParseWithClaims(tokenString, &dflclaims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, cher.New("unexpected_signing_method", cher.M{
						"alg": token.Header["alg"],
					})
				}

				return publicKey, nil
			})
			if err != nil {
				rpc.HandleError(w, r, err, nil)
				return
			}

			if !token.Valid {
				rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil), nil)
				return
			}

			claims := token.Claims.(*dfljwt.DFLClaims)

			authUser := authlib.AuthUser{
				ID:       claims.Subject,
				Username: claims.Username,
				Scopes:   claims.Scope,
			}

			ctx := context.WithValue(r.Context(), authlib.UserContextKey, authUser)

			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
