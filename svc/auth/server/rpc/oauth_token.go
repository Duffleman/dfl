package rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"dfl/lib/ptr"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"

	"github.com/xeipuuv/gojsonschema"
)

var tokenSchema = gojsonschema.NewStringLoader(`{
	"type": "object",
	"additionalProperties": false,

	"required": [
		"client_id",
		"grant_type",
		"redirect_uri",
		"code",
		"code_verifier"
	],

	"properties": {
		"client_id": {
			"type": "string",
			"minLength": 1
		},

		"grant_type": {
			"type": "string",
			"enum": ["authorization_code"]
		},

		"redirect_uri": {
			"type": ["string", "null"],
			"minLength": 1
		},

		"code": {
			"type": "string",
			"minLength": 1
		},

		"code_verifier": {
			"type": "string",
			"pattern": "^[A-Za-z\\d\\-\\._~]{43,128}$"
		}
	}
}`)

func Token(a *app.App, w http.ResponseWriter, r *http.Request) error {
	if err := modifyBody(r); err != nil {
		return err
	}

	if err := rpc.ValidateRequest(r, tokenSchema); err != nil {
		return err
	}

	req := &auth.TokenRequest{}
	if err := rpc.ParseBody(r, req); err != nil {
		return err
	}

	res, err := a.Token(r.Context(), req)
	if err != nil {
		return err
	}

	return rpc.WriteOut(w, res)
}

func modifyBody(r *http.Request) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}

	var redirectURI *string

	if v, ok := vals["redirect_uri"]; ok {
		redirectURI = ptr.String(v[0])
	}

	req := &auth.TokenRequest{
		ClientID:     vals.Get("client_id"),
		GrantType:    vals.Get("grant_type"),
		RedirectURI:  redirectURI,
		Code:         vals.Get("code"),
		CodeVerifier: vals.Get("code_verifier"),
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(jsonBytes))

	return nil
}
