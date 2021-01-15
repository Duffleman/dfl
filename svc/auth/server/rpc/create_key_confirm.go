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

func (r *RPC) CreateKeyConfirm(ctx context.Context, req *auth.CreateKeyConfirmRequest) error {
	user, err := r.app.FindUser(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := r.app.CheckLoginValidity(ctx, user); err != nil {
		return err
	}

	waUser, err := r.app.ConvertUserForWA(ctx, user, true)
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

	parsed, err := parseCredentialCreateKey(req)
	if err != nil {
		return err
	}

	credential, err := r.app.WA.CreateCredential(waUser, *session, parsed)
	if err != nil {
		return err
	}

	if _, err := r.app.CreateU2FCredential(ctx, user.ID, req.ChallengeID, req.KeyName, credential, nil); err != nil {
		return err
	}

	return nil
}

func parseCredentialCreateKey(req *auth.CreateKeyConfirmRequest) (*protocol.ParsedCredentialCreationData, error) {
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
