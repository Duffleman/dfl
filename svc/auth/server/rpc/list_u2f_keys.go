package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var listU2FKeysSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"include_unsigned"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		},

		"include_unsigned": {
			"type": "boolean"
		}
	}
}`)

func (r *RPC) ListU2FKeys(ctx context.Context, req *auth.ListU2FKeysRequest) ([]*auth.PublicU2FKey, error) {
	authUser := authlib.GetUserContext(ctx)

	if authUser.ID != req.UserID && !authUser.Can("auth:list_keys") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return r.app.ListU2FKeys(ctx, req)
}
