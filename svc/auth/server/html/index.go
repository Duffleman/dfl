package html

import (
	"net/http"

	"dfl/svc/auth/server/app"
)

func Index(a *app.App, w http.ResponseWriter, r *http.Request) error {
	return a.Template.Display(w, "index.html", map[string]interface{}{
		"title":      "DFL Auth",
		"activeHome": true,
	})
}
