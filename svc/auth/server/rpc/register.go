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

func Register(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rpc.ValidateRequest(r, registerSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &auth.RegisterRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		err = a.Register(r.Context(), req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		return
	}
}
