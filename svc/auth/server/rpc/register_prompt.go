package rpc

import (
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/xeipuuv/gojsonschema"
)

var registerSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username",
		"invite_code"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		},

		"invite_code": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func RegisterPrompt(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, registerSchema); err != nil {
		return err
	}

	req := &auth.RegisterPromptRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user, err := a.GetUserByName(r.Context(), req.Username)
	if err != nil {
		return err
	}

	if err := a.CheckRegistrationValidity(r.Context(), user, &req.InviteCode); err != nil {
		return err
	}

	waUser, err := a.ConvertUserForWA(r.Context(), user, true)
	if err != nil {
		return err
	}

	options, session, err := a.WA.BeginRegistration(waUser)

	for _, key := range waUser.Credentials {
		options.Response.CredentialExcludeList = append(options.Response.CredentialExcludeList, protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: key.ID,
		})
	}

	id, err := a.CreateU2FChallenge(r.Context(), session)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, &auth.RegisterPromptResponse{
		ID:        id,
		Challenge: options,
	})
}
