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

var createSignedURLSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"content_type"
	],

	"properties": {
		"name": {
			"type": "string",
			"minLength": 1
		},

		"content_type": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

// CreateSignedURL creates a signed URL for file uploads
func CreateSignedURL(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := rpc.ValidateRequest(r, createSignedURLSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &short.CreateSignedURLRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
		if !authUser.Can("short:upload") && !authUser.Can("short:admin") {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		res, err := a.CreateSignedURL(ctx, authUser.Username, req.Name, req.ContentType)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		accept := r.Header.Get("Accept")

		if strings.Contains(accept, "text/plain") {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(res.URL))
		} else {
			rpc.WriteOut(w, res)
		}
	}
}
