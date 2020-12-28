package rpc

import (
	"net/http"

	"dfl/svc/auth"
	"dfl/svc/auth/server/app"
	"dfl/svc/auth/server/lib/rpc"

	"github.com/xeipuuv/gojsonschema"
)

var loginSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username",
		"password"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		},

		"password": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func Login(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rpc.ValidateRequest(r, loginSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &auth.LoginRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		res, err := a.Login(r.Context(), req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		rpc.WriteOut(w, r, res)
		return
	}
}
