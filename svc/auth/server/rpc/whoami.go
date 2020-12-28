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

func WhoAmI(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rpc.ValidateRequest(r, whoamiSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &auth.WhoAmIRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		user := authlib.GetFromContext(r.Context())
		if user.Username != req.Username {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		res, err := a.WhoAmI(r.Context(), req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		rpc.WriteOut(w, r, res)
		return
	}
}
