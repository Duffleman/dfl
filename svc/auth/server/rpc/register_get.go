package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

func RegisterGet(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return rpc.QuickTemplate(w, map[string]interface{}{
		"title": "Register",
	}, []string{
		"./resources/register.html",
		"./resources/layouts/root.html",
	})
}
