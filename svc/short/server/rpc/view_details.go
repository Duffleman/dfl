package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/short"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var viewDetailsSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"query"
	],

	"properties": {
		"query": {
			"type": "string",
			"minLength": 3
		}
	}
}`)

func (r *RPC) ViewDetails(ctx context.Context, req *short.IdentifyResource) (*short.Resource, error) {
	authUser := authlib.GetUserContext(ctx)
	if !authUser.Can("short:upload") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	qi := r.app.ParseQueryType(req.Query)

	if len(qi) != 1 {
		return nil, cher.New("multi_query_not_supported", cher.M{"query": qi})
	}

	resource, err := r.app.GetResource(ctx, qi[0])
	if err != nil {
		return nil, err
	}

	if resource.Owner != authUser.Username && !authUser.Can("short:admin") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return resource, nil
}
