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
func CreateSignedURL(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	err := rpc.ValidateRequest(r, createSignedURLSchema)
	if err != nil {
		return err
	}

	req := &short.CreateSignedURLRequest{}
	err = rpc.ParseBody(r, req)
	if err != nil {
		return err
	}

	authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
	if !authUser.Can("short:upload") && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	res, err := a.CreateSignedURL(ctx, authUser.Username, req.Name, req.ContentType)
	if err != nil {
		return err
	}

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/plain") {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(res.URL))
	} else {
		return rpc.WriteOut(w, res)
	}

	return nil
}
