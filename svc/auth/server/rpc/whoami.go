package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var whoamiSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func WhoAmI(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, whoamiSchema); err != nil {
		return err
	}

	req := &auth.WhoAmIRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user := authlib.GetFromContext(r.Context())
	if user.Username != req.Username {
		return cher.New(cher.AccessDenied, nil)
	}

	if !user.Can("auth:whoami") {
		return cher.New(cher.AccessDenied, nil)
	}

	res, err := a.WhoAmI(r.Context(), req)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, r, res)
}
