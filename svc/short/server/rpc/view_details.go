package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/short"
	"dfl/svc/short/server/app"

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

func ViewDetails(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	err := rpc.ValidateRequest(r, viewDetailsSchema)
	if err != nil {
		return err
	}

	req := &short.IdentifyResource{}
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

	resource, err := a.GetResource(ctx, qi[0])
	if err != nil {
		return err
	}

	if resource.Owner != authUser.Username && !authUser.Can("short:admin") {
		return cher.New(cher.AccessDenied, nil)
	}

	return rpc.WriteOut(w, resource)
}
