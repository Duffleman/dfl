package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var listU2FKeysSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"include_unsigned"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		},

		"include_unsigned": {
			"type": "boolean"
		}
	}
}`)

func ListU2FKeys(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, listU2FKeysSchema); err != nil {
		return err
	}

	req := &auth.ListU2FKeysRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user := authlib.GetFromContext(r.Context())
	if user.ID != req.UserID && !user.Can("auth:list_keys") {
		return cher.New(cher.AccessDenied, nil)
	}

	res, err := a.ListU2FKeys(r.Context(), req.UserID, req.IncludeUnsigned)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, res)
}
