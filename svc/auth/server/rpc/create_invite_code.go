package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/xeipuuv/gojsonschema"
)

var createInviteCodeSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"scopes",
		"code",
		"expires_at"
	],

	"properties": {
		"scopes": {
			"type": "string",
			"minLength": 1
		},

		"code": {
			"type": ["string", "null"],
			"minLength": 1
		},

		"expires_at": {
			"type": ["string", "null"],
			"format": "date-time"
		}
	}
}`)

func CreateInviteCode(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, createInviteCodeSchema); err != nil {
		return err
	}

	req := &auth.CreateInviteCodeRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	authUser := authlib.GetFromContext(r.Context())
	if !authUser.Can("auth:create_invite_code") {
		return cher.New(cher.AccessDenied, nil)
	}

	user, err := a.FindUser(r.Context(), authUser.ID)
	if err != nil {
		return err
	}

	if !authlib.Can(req.Scopes, user.Scopes) {
		return cher.New(cher.AccessDenied, nil)
	}

	res, err := a.CreateInviteCode(r.Context(), authUser.ID, req)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, res)
}
