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

	if len(templates) == 1 {
		return tpl.ExecuteTemplate(w, firstName, data)
	}

	return tpl.ExecuteTemplate(w, "root", data)
}
