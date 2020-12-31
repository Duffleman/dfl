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

var createClientSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"name",
		"redirect_uris"
	],

	"properties": {
		"name": {
			"type": "string",
			"minLength": 3
		},

		"redirect_uris": {
			"type": "array",
			"minItems": 0,

			"items": {
				"type": "string",
				"minLength": 1
			}
		}
	}
}`)

func CreateClient(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, createClientSchema); err != nil {
		return err
	}

	req := &auth.CreateClientRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user := authlib.GetFromContext(r.Context())
	if !user.Can("auth:create_client") {
		return cher.New(cher.AccessDenied, nil)
	}

	res, err := a.CreateClient(r.Context(), req)
	if err != nil {
		return err
	}

	rpc.WriteOut(w, r, res)
	return nil
}
