package rpc

import (
	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"
	"net/http"

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

func DeleteKey(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, deleteKeySchema); err != nil {
		return err
	}

	req := &auth.DeleteKeyRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user := authlib.GetFromContext(r.Context())
	if user.ID != req.UserID && !user.Can("auth:delete_keys") {
		return cher.New(cher.AccessDenied, nil)
	}

	return a.DeleteKey(r.Context(), req.UserID, req.KeyID)
}
