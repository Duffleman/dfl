package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var shortenURLSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"url"
	],

	"properties": {
		"url": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func (r *RPC) ShortenURL(ctx context.Context, req *short.CreateURLRequest) (*short.CreateResourceResponse, error) {
	authUser := authlib.GetUserContext(ctx)
	if !authUser.Can("short:upload") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return r.app.ShortenURL(ctx, req.URL, authUser.Username)
}
