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

		user, err := a.GetUserByName(r.Context(), req.Username)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if !authlib.Can("auth:login", user.Scopes) {
			rpc.HandleError(w, r, cher.New(cher.Unauthorized, nil))
			return
		}

		res, err := a.Login(r.Context(), req, user)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		rpc.WriteOut(w, r, res)
		return
	}
}
