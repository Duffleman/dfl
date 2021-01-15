package rpc

import (
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/xeipuuv/gojsonschema"
)

var createKeyPromptSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		}
	}
}`)

func CreateKeyPrompt(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, createKeyPromptSchema); err != nil {
		return err
	}

	req := &auth.CreateKeyPromptRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	authuser := authlib.GetFromContext(r.Context())
	if authuser.ID != req.UserID && !authuser.Can("auth:create_keys") {
		return cher.New(cher.AccessDenied, nil)
	}

	user, err := a.FindUser(r.Context(), req.UserID)
	if err != nil {
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

	return rpc.WriteOut(w, &auth.CreateKeyPromptResponse{
		ID:        id,
		Challenge: options,
	})
}
