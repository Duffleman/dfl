package templates

import (
	"embed"
	"net/http"
	"text/template"
)

type Template struct {
	*template.Template
}

func New(fs embed.FS) (*Template, error) {
	tpl, err := template.New("").Delims("[[", "]]").ParseFS(fs, "resources/*")
	if err != nil {
		return nil, err
	}

	return &Template{tpl}, nil
}

func (t *Template) Display(w http.ResponseWriter, pageName string, data map[string]interface{}) error {
	return t.ExecuteTemplate(w, pageName, data)
}
