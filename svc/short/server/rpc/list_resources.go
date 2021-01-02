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

func ListResources(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	err := rpc.ValidateRequest(r, listResourcesSchema)
	if err != nil {
		return err
	}

	req := &short.ListResourcesRequest{}
	err = rpc.ParseBody(r, req)
	if err != nil {
		return err
	}

	authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
	if !authUser.Can("short:upload") && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	if err = authorizeRequest(req, authUser); err != nil {
		return err
	}

	resources, err := a.ListResources(ctx, req)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, resources)
}

func authorizeRequest(req *short.ListResourcesRequest, u authlib.AuthUser) error {
	switch {
	case u.Can("short:admin"):
		return nil
	case req.Username != nil && *req.Username == u.Username:
		return nil
	default:
		return cher.New(cher.AccessDenied, nil)
	}
}
