package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var deleteKeySchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"key_id"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		},

		"key_id": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func (r *RPC) DeleteKey(ctx context.Context, req *auth.DeleteKeyRequest) error {
	authUser := authlib.GetUserContext(ctx)

	if authUser.ID != req.UserID && !authUser.Can("auth:delete_keys") {
		return cher.New(cher.AccessDenied, nil)
	}

	return r.app.DeleteKey(ctx, req.UserID, req.KeyID)
}
