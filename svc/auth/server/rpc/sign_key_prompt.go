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

var signKeyPromptSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"key_to_sign"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		},

		"key_to_sign": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func SignKeyPrompt(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, signKeyPromptSchema); err != nil {
		return err
	}

	req := &auth.SignKeyPromptRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	authuser := authlib.GetFromContext(r.Context())
	if authuser.ID != req.UserID {
		return cher.New(cher.AccessDenied, nil)
	}

	user, err := a.FindUser(r.Context(), req.UserID)
	if err != nil {
		return err
	}

	if err := a.CanSign(r.Context(), user.ID, req.KeyToSign); err != nil {
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

	return rpc.WriteOut(w, &auth.SignKeyPromptResponse{
		ID:        id,
		Challenge: options,
	})
}
