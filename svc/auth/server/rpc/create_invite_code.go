package rpc

import (
	"context"

	authlib "dfl/lib/auth"
	"dfl/svc/auth"

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

func (r *RPC) CreateInviteCode(ctx context.Context, req *auth.CreateInviteCodeRequest) (*auth.CreateInviteCodeResponse, error) {
	authUser := authlib.GetUserContext(ctx)

	if !authUser.Can("auth:create_invite_code") {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	user, err := r.app.FindUser(ctx, authUser.ID)
	if err != nil {
		return nil, err
	}

	if !authlib.Can(req.Scopes, user.Scopes) {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	return r.app.CreateInviteCode(ctx, user.ID, req)
}
