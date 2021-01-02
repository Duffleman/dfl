package rpc

import (
	"net/http"
	"path"
	"text/template"
)

func QuickTemplate(w http.ResponseWriter, data interface{}, templates []string) error {
	_, firstName := path.Split(templates[0])

	tpl, err := template.New(firstName).Delims("[[", "]]").ParseFiles(templates...)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(w, "root", data)
}
