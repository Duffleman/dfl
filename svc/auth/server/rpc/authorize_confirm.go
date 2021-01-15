package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	"dfl/svc/auth"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/xeipuuv/gojsonschema"
)

var authorizeConfirmSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"response_type",
		"redirect_uri",
		"client_id",
		"scope",
		"state",
		"nonce",
		"code_challenge",
		"code_challenge_method",
		"username",
		"challenge_id",
		"webauthn"
	],

	"properties": {
		"response_type": {
			"type": "string",
			"enum": ["code"]
		},

		"redirect_uri": {
			"type": ["string", "null"],
			"minLength": 1
		},

		"client_id": {
			"type": "string",
			"minLength": 1
		},

		"scope": {
			"type": "string",
			"minLength": 1
		},

		"state": {
			"type": "string",
			"minLength": 1
		},

		"nonce": {
			"type": "string",
			"minLength": 1
		},

		"code_challenge": {
			"type": "string",
			"minLength": 1
		},

		"code_challenge_method": {
			"type": "string",
			"enum": ["S256"]
		},

		"username": {
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

func (r *RPC) AuthorizeConfirm(ctx context.Context, req *auth.AuthorizeConfirmRequest) (*auth.AuthorizeConfirmResponse, error) {
	user, err := r.app.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	waUser, err := r.app.ConvertUserForWA(ctx, user, false)
	if err != nil {
		return nil, err
	}

	session, err := r.app.FindU2FChallenge(ctx, req.ChallengeID)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(session.UserID, waUser.WebAuthnID()) != 0 {
		return nil, cher.New(cher.NotFound, nil) // pretend to not know whats going on
	}

	parsed, err := parseCredentialLogin(req)
	if err != nil {
		return nil, err
	}

	if _, err := r.app.WA.ValidateLogin(waUser, *session, parsed); err != nil {
		return nil, err
	}

	return r.app.Authorization(ctx, req, user)
}

func parseCredentialLogin(req *auth.AuthorizeConfirmRequest) (*protocol.ParsedCredentialAssertionData, error) {
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
