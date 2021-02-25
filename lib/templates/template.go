package templates

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

type Template struct {
	fs   embed.FS
	base *template.Template
}

func New(fs embed.FS) *Template {
	tpl := template.New("").Delims("[[", "]]")

	return &Template{fs, tpl}
}

func (t *Template) Display(w http.ResponseWriter, pageName string, data map[string]interface{}) error {
	pagePath := fmt.Sprintf("resources/%s.html", pageName)

	tpl, err := t.base.ParseFS(t.fs, "resources/root.html", "resources/_nav.html", pagePath)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(w, "root", data)
}
