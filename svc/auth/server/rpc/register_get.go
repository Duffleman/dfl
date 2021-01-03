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
		"./resources/auth/register.html",
		"./resources/auth/layouts/root.html",
	})
}
