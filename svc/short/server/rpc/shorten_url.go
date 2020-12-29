package rpc

import (
	"net/http"
	"strings"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/short"
	"dfl/svc/short/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var shortenURLSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"url"
	],

	"properties": {
		"url": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func ShortenURL(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := rpc.ValidateRequest(r, shortenURLSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &short.CreateURLRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
		if !authUser.Can("short:upload") {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		res, err := a.ShortenURL(ctx, req.URL, authUser.Username)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		accept := r.Header.Get("Accept")

		if strings.Contains(accept, "text/plain") {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(res.URL))
		} else {
			rpc.WriteOut(w, r, res)
		}
	}
}
