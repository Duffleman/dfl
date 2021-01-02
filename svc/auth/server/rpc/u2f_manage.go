package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

func U2FManageGet(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return rpc.QuickTemplate(w, map[string]interface{}{
		"title": "Manage U2F",
	}, []string{
		"./resources/u2f_manage.html",
		"./resources/layouts/root.html",
	})
}
