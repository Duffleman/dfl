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

var deleteResourceSchema = gojsonschema.NewStringLoader(`{
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

func DeleteResource(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := rpc.ValidateRequest(r, deleteResourceSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &short.IdentifyResource{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		authUser := ctx.Value(authlib.UserContextKey).(authlib.AuthUser)
		if !authUser.Can("short:delete") {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		qi := a.ParseQueryType(req.Query)

		if len(qi) != 1 {
			rpc.HandleError(w, r, cher.New("multi_query_not_supported", cher.M{"query": qi}))
			return
		}

		if qi[0].QueryType == app.Name {
			rpc.HandleError(w, r, cher.New("cannot_query_resource_by_name", cher.M{"query": qi}))
			return
		}

		resource, err := a.GetResource(ctx, qi[0])
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if resource.Owner != authUser.Username {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil))
			return
		}

		err = a.DeleteResource(ctx, resource)
		rpc.HandleError(w, r, err)
	}
}
