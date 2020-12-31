package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var registerSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username",
		"password",
		"invite_code"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		},

		"email": {
			"type": "string",
			"format": "email"
		},

		"password": {
			"type": "string",
			"minLength": 1
		},

		"invite_code": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func Register(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, registerSchema); err != nil {
		return err
	}

	req := &auth.RegisterRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	return a.Register(r.Context(), req)
}
