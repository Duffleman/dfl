package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

var registerPage = rpc.MakeTemplate([]string{
	"./resources/auth/register.html",
	"./resources/auth/_nav.html",
	"./resources/auth/layouts/root.html",
})

func RegisterGet(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return registerPage.Execute(w, map[string]interface{}{
		"title":          "Register",
		"activeRegister": true,
	})
}
