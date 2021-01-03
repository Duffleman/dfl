package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

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

func SignKeyConfirm(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, signKeyConfirmSchema); err != nil {
		return err
	}

	req := &auth.SignKeyConfirmRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	authUser := authlib.GetFromContext(r.Context())
	if authUser.ID != req.UserID {
		return cher.New(cher.AccessDenied, nil)
	}

	user, err := a.FindUser(r.Context(), req.UserID)
	if err != nil {
		return err
	}

	waUser, err := a.ConvertUserForWA(r.Context(), user, false)
	if err != nil {
		return err
	}

	session, err := a.FindU2FChallenge(r.Context(), req.ChallengeID)
	if err != nil {
		return err
	}

	if bytes.Compare(session.UserID, waUser.WebAuthnID()) != 0 {
		return cher.New(cher.NotFound, nil) // pretend to not know whats going on
	}

	parsed, err := ParseCredentialSignKey(req)
	if err != nil {
		return err
	}

	if _, err := a.WA.ValidateLogin(waUser, *session, parsed); err != nil {
		return err
	}

	return a.SignKey(r.Context(), user, req.KeyToSign)
}

func ParseCredentialSignKey(req *auth.SignKeyConfirmRequest) (*protocol.ParsedCredentialAssertionData, error) {
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
