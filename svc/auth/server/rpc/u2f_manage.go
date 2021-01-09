package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

var managePage = rpc.MakeTemplate([]string{
	"./resources/auth/u2f_manage.html",
	"./resources/auth/_nav.html",
	"./resources/auth/layouts/root.html",
})

func U2FManageGet(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return managePage.Execute(w, map[string]interface{}{
		"title":           "Manage U2F",
		"activeManageU2F": true,
	})
}
