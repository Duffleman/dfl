package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var tokenSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"client_id",
		"grant_type",
		"code",
		"code_verifier"
	],

	"properties": {
		"client_id": {
			"type": "string",
			"minLength": 1
		},

		"grant_type": {
			"type": "string",
			"enum": ["code"]
		},

		"redirect_uri": {
			"type": "string",
			"minLength": 1
		},

		"code": {
			"type": "string",
			"minLength": 1
		},

		"code_verifier": {
			"type": "string",
			"pattern": "^[A-Za-z\\d\\-\\._~]{43,128}$"
		}
	}
}`)

func Token(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rpc.ValidateRequest(r, tokenSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &auth.TokenRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		res, err := a.Token(r.Context(), req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		rpc.WriteOut(w, r, res)
		return
	}
}
