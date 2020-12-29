package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/short/server/app"
)

func ResaveHashes(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := a.ResaveHashes(ctx)
		rpc.HandleError(w, r, err)
	}
}
