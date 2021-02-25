package html

import (
	"net/http"

	"dfl/svc/auth/server/app"
)

func U2FManage(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return a.Template.Display(w, "u2f_manage", map[string]interface{}{
		"title":           "Manage U2F",
		"activeManageU2F": true,
	})
}
