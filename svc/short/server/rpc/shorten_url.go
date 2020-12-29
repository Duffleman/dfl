package rpc

import (
	"net/http"
	"strings"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/short"
	"dfl/svc/short/server/app"

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

func ShortenURL(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := rpc.ValidateRequest(r, shortenURLSchema)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		req := &short.CreateURLRequest{}
		err = rpc.ParseBody(r, req)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		username := ctx.Value(authlib.UserContextKey).(string)

		res, err := a.ShortenURL(ctx, req.URL, username)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		accept := r.Header.Get("Accept")

		if strings.Contains(accept, "text/plain") {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(res.URL))
		} else {
			rpc.WriteOut(w, r, res)
		}
	}
}
