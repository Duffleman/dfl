package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var listResourcesSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"include_deleted"
	],

	"properties": {
		"include_deleted": {
			"type": "boolean"
		},

		"username": {
			"type": "string",
			"minLength": 1
		},

		"limit": {
			"type": "number",
			"minimum": 1,
			"maximum": 100
		},

		"filter_mime": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func (r *RPC) ListResources(ctx context.Context, req *short.ListResourcesRequest) ([]*short.Resource, error) {
	authUser := authlib.GetUserContext(ctx)

	if err := authorizeRequest(req, authUser); err != nil {
		return nil, err
	}

	return r.app.ListResources(ctx, req)
}

func authorizeRequest(req *short.ListResourcesRequest, u *authlib.AuthUser) error {
	switch {
	case u.Can("short:admin"):
		return nil
	case req.Username != nil && *req.Username == u.Username:
		return nil
	default:
		return cher.New(cher.AccessDenied, nil)
	}
}
