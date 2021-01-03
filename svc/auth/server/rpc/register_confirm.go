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

var registerConfirmSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"username",
		"invite_code",
		"challenge_id",
		"webauthn"
	],

	"properties": {
		"username": {
			"type": "string",
			"minLength": 1
		},

		"invite_code": {
			"type": "string",
			"minLength": 1
		},

		"key_name": {
			"type": ["null", "string"],
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
						"attestation_object",
						"client_data_json"
					],

					"properties": {
						"attestation_object": {
							"type": "string",
							"minLength": 1
						},

						"client_data_json": {
							"type": "string",
							"minLength": 1
						}
					}
				}
			}
		}
	}
}`)

func RegisterConfirm(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, registerConfirmSchema); err != nil {
		return err
	}

	req := &auth.RegisterConfirmRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	session, err := a.FindU2FChallenge(r.Context(), req.ChallengeID)
	if err != nil {
		return err
	}

	user := &auth.User{
		ID:       string(session.UserID),
		Username: req.Username,
	}

	waUser, err := a.ConvertUserForWA(r.Context(), user, true)
	if err != nil {
		return err
	}

	if bytes.Compare(session.UserID, waUser.WebAuthnID()) != 0 {
		return cher.New(cher.NotFound, nil) // pretend to not know whats going on
	}

	parsed, err := ParseCredentialRegister(req)
	if err != nil {
		return err
	}

	credential, err := a.WA.CreateCredential(waUser, *session, parsed)
	if err != nil {
		return err
	}

	if _, err := a.Register(r.Context(), req, credential); err != nil {
		return err
	}

	return nil
}

func ParseCredentialRegister(req *auth.RegisterConfirmRequest) (*protocol.ParsedCredentialCreationData, error) {
	mm := map[string]interface{}{
		"id":    req.WebAuthn.ID,
		"rawId": req.WebAuthn.RawID,
		"type":  req.WebAuthn.Type,
		"response": map[string]interface{}{
			"attestationObject": req.WebAuthn.Response.AttestationObject,
			"clientDataJSON":    req.WebAuthn.Response.ClientDataJSON,
		},
	}

	rawBytes, err := json.Marshal(mm)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(rawBytes)

	parsed, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
