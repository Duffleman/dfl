package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/svc/short/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

func ResaveHashes(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
	if !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	return a.ResaveHashes(ctx)
}
