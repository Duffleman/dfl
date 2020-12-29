package rpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"dfl/lib/rpc"
)

func Robots() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileContent, err := ioutil.ReadFile("./resources/robots.txt")
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		modTime, err := time.Parse(time.RFC3339, "2019-10-02T12:00:00Z")
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		reader := bytes.NewReader(fileContent)

		http.ServeContent(w, r, "robots.txt", modTime, reader)
	}
}
