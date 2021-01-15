package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"

	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/xeipuuv/gojsonschema"
)

var createKeyConfirmSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"user_id",
		"challenge_id",
		"webauthn"
	],

	"properties": {
		"user_id": {
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

func CreateKeyConfirm(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := rpc.ValidateRequest(r, createKeyConfirmSchema); err != nil {
		return err
	}

	req := &auth.CreateKeyConfirmRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	user, err := a.FindUser(r.Context(), req.UserID)
	if err != nil {
		return err
	}

	if err := a.CheckLoginValidity(r.Context(), user); err != nil {
		return err
	}

	waUser, err := a.ConvertUserForWA(r.Context(), user, true)
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

	parsed, err := ParseCredentialCreateKey(req)
	if err != nil {
		return err
	}

	credential, err := a.WA.CreateCredential(waUser, *session, parsed)
	if err != nil {
		return err
	}

	if _, err := a.CreateU2FCredential(r.Context(), user.ID, req.ChallengeID, req.KeyName, credential, nil); err != nil {
		return err
	}

	return nil
}

func ParseCredentialCreateKey(req *auth.CreateKeyConfirmRequest) (*protocol.ParsedCredentialCreationData, error) {
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
