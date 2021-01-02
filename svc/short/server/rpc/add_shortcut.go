package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/short"
	"dfl/svc/short/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var addShortcutSchema = gojsonschema.NewStringLoader(`{
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

func AddShortcut(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	err := rpc.ValidateRequest(r, addShortcutSchema)
	if err != nil {
		return err
	}

	req := &short.ChangeShortcutRequest{}
	err = rpc.ParseBody(r, req)
	if err != nil {
		return err
	}

	authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
	if !authUser.Can("short:upload") && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	qi := a.ParseQueryType(req.Query)

	if len(qi) != 1 {
		return cher.New("multi_query_not_supported", cher.M{"query": qi})
	}

	if qi[0].QueryType == app.Name {
		return cher.New("cannot_query_resource_by_name", cher.M{"query": qi})
	}

	resource, err := a.GetResource(ctx, qi[0])
	if err != nil {
		return err
	}

	if resource.Owner != authUser.Username && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	return a.AddShortcut(ctx, resource, req.Shortcut)
}
