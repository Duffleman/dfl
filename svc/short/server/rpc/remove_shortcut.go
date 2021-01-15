package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/short"
	"dfl/svc/short/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var removeShortcutSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"query",
		"shortcut"
	],

	"properties": {
		"query": {
			"type": "string",
			"minLength": 3
		},

		"shortcut": {
			"type": "string",
			"minLength": 3
		}
	}
}`)

func (r *RPC) RemoveShortcut(ctx context.Context, req *short.ChangeShortcutRequest) error {
	authUser := authlib.GetUserContext(ctx)
	if !authUser.Can("short:upload") {
		return cher.New(cher.AccessDenied, nil)
	}

	qi := r.app.ParseQueryType(req.Query)

	if len(qi) != 1 {
		return cher.New("multi_query_not_supported", cher.M{"query": qi})
	}

	if qi[0].QueryType == app.Name {
		return cher.New("cannot_query_resource_by_name", cher.M{"query": qi})
	}

	resource, err := r.app.GetResource(ctx, qi[0])
	if err != nil {
		return err
	}

	if resource.Owner != authUser.Username && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	return r.app.RemoveShortcut(ctx, resource, req.Shortcut)
}
