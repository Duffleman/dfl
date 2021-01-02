package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

func WhoAmI(a *app.App, w http.ResponseWriter, r *http.Request) error {
	user := authlib.GetFromContext(r.Context())

	res, err := a.WhoAmI(r.Context(), user.ID)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, res)
}
