package rpc

import (
	"net/http"
	"path"
	"strings"
	"text/template"
)

type HTMLPage struct {
	Name     string
	Template *template.Template
}

func (p HTMLPage) Execute(w http.ResponseWriter, data interface{}) error {
	return p.Template.ExecuteTemplate(w, p.Name, data)
}

func MakeTemplate(templates []string) *HTMLPage {
	_, firstName := path.Split(templates[0])

	tpl, err := template.New(firstName).Delims("[[", "]]").ParseFiles(templates...)
	if err != nil {
		panic(err)
	}

	if len(templates) == 1 {
		return &HTMLPage{
			Name:     firstName,
			Template: tpl,
		}
	}

	lastItem := templates[len(templates)-1]
	_, file := path.Split(lastItem)
	ext := path.Ext(file)
	file = strings.TrimSuffix(file, ext)

	return &HTMLPage{
		Name:     file,
		Template: tpl,
	}
}
