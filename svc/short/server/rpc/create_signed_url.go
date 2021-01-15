package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var createSignedURLSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"content_type"
	],

	"properties": {
		"name": {
			"type": "string",
			"minLength": 1
		},

		"content_type": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

// CreateSignedURL creates a signed URL for file uploads
func (r *RPC) CreateSignedURL(ctx context.Context, req *short.CreateSignedURLRequest) (*short.CreateSignedURLResponse, error) {
	authUser := authlib.GetUserContext(ctx)
	if !authUser.Can("short:upload") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return r.app.CreateSignedURL(ctx, authUser.Username, req.Name, req.ContentType)
}
