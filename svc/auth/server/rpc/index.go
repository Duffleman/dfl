package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

var indexPage = rpc.MakeTemplate([]string{
	"./resources/auth/index.html",
	"./resources/auth/_nav.html",
	"./resources/auth/layouts/root.html",
})

func Index(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return indexPage.Execute(w, map[string]interface{}{
		"title":      "DFL Auth",
		"activeHome": true,
	})
}
