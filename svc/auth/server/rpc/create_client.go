package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
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

func (r *RPC) CreateClient(ctx context.Context, req *auth.CreateClientRequest) (*auth.CreateClientResponse, error) {
	authUser := authlib.GetUserContext(ctx)

	if !authUser.Can("auth:create_client") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return r.app.CreateClient(ctx, req)
}
