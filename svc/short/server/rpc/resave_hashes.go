package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/short/server/app"
)

func ResaveHashes(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
		if !authUser.Can("short:admin") {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		err := a.ResaveHashes(ctx)
		rpc.HandleError(w, r, err)
	}
}
