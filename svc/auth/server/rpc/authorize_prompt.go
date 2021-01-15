package rpc

import (
	"context"

	"dfl/svc/auth"

	"github.com/xeipuuv/gojsonschema"
)

var authorizePromptSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func (r *RPC) AuthorizePrompt(ctx context.Context, req *auth.AuthorizePromptRequest) (*auth.AuthorizePromptResponse, error) {
	user, err := r.app.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := r.app.CheckLoginValidity(ctx, user); err != nil {
		return nil, err
	}

	waUser, err := r.app.ConvertUserForWA(ctx, user, false)
	if err != nil {
		return nil, err
	}

	options, session, err := r.app.WA.BeginLogin(waUser)
	if err != nil {
		return nil, err
	}

	id, err := r.app.CreateU2FChallenge(ctx, session)
	if err != nil {
		return nil, err
	}

	return &auth.AuthorizePromptResponse{
		ID:        id,
		Challenge: options,
	}, nil
}
