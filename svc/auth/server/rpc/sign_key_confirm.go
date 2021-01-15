package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	authlib "dfl/lib/auth"
	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/xeipuuv/gojsonschema"
)

var signKeyConfirmSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"key_to_sign",
		"challenge_id",
		"webauthn"
	],

	"properties": {
		"user_id": {
			"type": "string",
			"minLength": 1
		},

		"key_to_sign": {
			"type": "string",
			"minLength": 1
		},

		"challenge_id": {
			"type": "string",
			"minLength": 1
		},

		"webauthn": {
			"type": "object",
			"additionalProperties": false,

			"required": [
				"id",
				"raw_id",
				"type",
				"response"
			],

			"properties": {
				"id": {
					"type": "string",
					"minLength": 1
				},

				"raw_id": {
					"type": "string",
					"minLength": 1
				},

				"type": {
					"type": "string",
					"minLength": 1
				},

				"response": {
					"type": "object",
					"additionalProperties": false,

					"required": [
						"authenticator_data",
						"client_data_json",
						"signature",
						"user_handle"
					],

					"properties": {
						"authenticator_data": {
							"type": "string",
							"minLength": 1
						},

						"client_data_json": {
							"type": "string",
							"minLength": 1
						},

						"signature": {
							"type": "string",
							"minLength": 1
						},

						"user_handle": {
							"type": "string"
						}
					}
				}
			}
		}
	}
}`)

func (r *RPC) SignKeyConfirm(ctx context.Context, req *auth.SignKeyConfirmRequest) error {
	authUser := authlib.GetUserContext(ctx)
	if authUser.ID != req.UserID {
		return cher.New(cher.AccessDenied, nil)
	}

	user, err := r.app.FindUser(ctx, req.UserID)
	if err != nil {
		return err
	}

	waUser, err := r.app.ConvertUserForWA(ctx, user, false)
	if err != nil {
		return err
	}

	session, err := r.app.FindU2FChallenge(ctx, req.ChallengeID)
	if err != nil {
		return err
	}

	if bytes.Compare(session.UserID, waUser.WebAuthnID()) != 0 {
		return cher.New(cher.NotFound, nil) // pretend to not know whats going on
	}

	parsed, err := parseCredentialSignKey(req)
	if err != nil {
		return err
	}

	if _, err := r.app.WA.ValidateLogin(waUser, *session, parsed); err != nil {
		return err
	}

	return r.app.SignKey(ctx, user, req.KeyToSign)
}

func parseCredentialSignKey(req *auth.SignKeyConfirmRequest) (*protocol.ParsedCredentialAssertionData, error) {
	mm := map[string]interface{}{
		"id":    req.WebAuthn.ID,
		"rawId": req.WebAuthn.RawID,
		"type":  req.WebAuthn.Type,
		"response": map[string]interface{}{
			"authenticatorData": req.WebAuthn.Response.AuthenticatorData,
			"clientDataJSON":    req.WebAuthn.Response.ClientDataJSON,
			"signature":         req.WebAuthn.Response.Signature,
			"userHandle":        req.WebAuthn.Response.UserHandle,
		},
	}

	rawBytes, err := json.Marshal(mm)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(rawBytes)

	parsed, err := protocol.ParseCredentialRequestResponseBody(reader)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
