package rpc

import (
	"net/http"

	"dfl/lib/cher"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"
	authlib "dfl/svc/auth/server/lib/auth"
	"dfl/svc/auth/server/lib/rpc"

	"github.com/xeipuuv/gojsonschema"
)

var createClientSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"name"
	],

	"properties": {
		"name": {
			"type": "string",
			"minLength": 3
		}
	}
}`)

func CreateClient(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rpc.ValidateRequest(r, createClientSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &auth.CreateClientRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		user := authlib.GetFromContext(r.Context())
		if !user.Can("dflauth:create_client") {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		res, err := a.CreateClient(r.Context(), req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		rpc.WriteOut(w, r, res)
		return
	}
}
