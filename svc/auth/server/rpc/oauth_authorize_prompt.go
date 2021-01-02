package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

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

func AuthorizePrompt(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, authorizePromptSchema); err != nil {
		return err
	}

	req := &auth.AuthorizePromptRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user, err := a.GetUserByName(r.Context(), req.Username)
	if err != nil {
		return err
	}

	if err := a.CheckLoginValidity(r.Context(), user); err != nil {
		return err
	}

	waUser, err := a.ConvertUserForWA(r.Context(), user, false)
	if err != nil {
		return err
	}

	options, session, err := a.WA.BeginLogin(waUser)
	if err != nil {
		return err
	}

	id, err := a.CreateU2FChallenge(r.Context(), session)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, &auth.AuthorizePromptResponse{
		ID:        id,
		Challenge: options,
	})
}
