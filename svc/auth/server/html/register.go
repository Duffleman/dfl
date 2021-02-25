package html

import (
	"net/http"

	"dfl/svc/auth/server/app"
)

func Register(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return a.Template.Display(w, "register", map[string]interface{}{
		"title":          "Register",
		"activeRegister": true,
	})
}
