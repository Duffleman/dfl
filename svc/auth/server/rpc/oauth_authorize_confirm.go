package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"

	"dfl/lib/cher"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

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

func AuthorizeConfirm(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, authorizeConfirmSchema); err != nil {
		return err
	}

	req := &auth.AuthorizeConfirmRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user, err := a.GetUserByName(r.Context(), req.Username)
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

	parsed, err := ParseCredentialLogin(req)
	if err != nil {
		return err
	}

	if _, err := a.WA.ValidateLogin(waUser, *session, parsed); err != nil {
		return err
	}

	res, err := a.Authorization(r.Context(), req, user)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, res)
}

func ParseCredentialLogin(req *auth.AuthorizeConfirmRequest) (*protocol.ParsedCredentialAssertionData, error) {
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
