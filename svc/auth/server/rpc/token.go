package rpc

import (
	"context"

	"dfl/svc/auth"

	"github.com/xeipuuv/gojsonschema"
)

var tokenSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"client_id",
		"grant_type",
		"redirect_uri",
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
			"enum": ["authorization_code"]
		},

		"redirect_uri": {
			"type": ["string", "null"],
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

func (r *RPC) Token(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	return r.app.Token(ctx, req)
}
