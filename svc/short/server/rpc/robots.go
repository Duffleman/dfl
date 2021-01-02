package rpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"dfl/svc/short/server/app"
)

func Robots(_ *app.App, w http.ResponseWriter, r *http.Request) error {
	fileContent, err := ioutil.ReadFile("./resources/robots.txt")
	if err != nil {
		return err
	}

	modTime, err := time.Parse(time.RFC3339, "2019-10-02T12:00:00Z")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(fileContent)

	http.ServeContent(w, r, "robots.txt", modTime, reader)

	return nil
}
