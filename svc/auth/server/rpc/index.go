package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth/server/app"
)

func Index(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return rpc.QuickTemplate(w, map[string]interface{}{
		"title": "DFL Auth",
	}, []string{
		"./resources/auth/index.html",
		"./resources/auth/layouts/root.html",
	})
}
