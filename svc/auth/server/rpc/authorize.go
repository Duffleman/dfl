package rpc

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"dfl/lib/cher"
	"dfl/lib/ptr"
	"dfl/lib/rpc"
	"dfl/svc/auth"
	"dfl/svc/auth/server/app"
)

func AuthorizeGet(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := parseAuthorizeParams(r)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		client, err := a.DB.Q.FindClient(r.Context(), params.ClientID)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		tpl, err := template.ParseFiles("./resources/auth.html")
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		err = tpl.Execute(w, map[string]interface{}{
			"client_name": client.Name,
			"params":      params,
			"scopes":      strings.Fields(params.Scope),
		})
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}
	}
}

func AuthorizePost(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authCredentials, params, err := parseAuthorizeResponse(r)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		client, err := a.DB.Q.FindClient(r.Context(), params.ClientID)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		res, err := a.Authorization(r.Context(), &auth.AuthorizationRequest{
			ClientID:            client.ID,
			ResponseType:        params.ResponseType,
			RedirectURI:         params.RedirectURI,
			State:               params.State,
			CodeChallengeMethod: params.CodeChallengeMethod,
			CodeChallenge:       params.CodeChallenge,
			Username:            authCredentials.Username,
			Password:            authCredentials.Password,
			Scope:               params.Scope,
		})
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		tpl, err := template.ParseFiles("./resources/auth_code.html")
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		t, _ := time.Parse(time.RFC3339, res.ExpiresAt)

		err = tpl.Execute(w, map[string]interface{}{
			"client_name":          client.Name,
			"code":                 res.AuthorizationCode,
			"expires_at":           res.ExpiresAt,
			"expires_in":           res.ExpiresIn,
			"expires_at_formatted": t.Format(time.RFC822),
			"state":                res.State,
		})

		rpc.HandleError(w, r, err)
	}
}

type params struct {
	ResponseType        string  `json:"response_type"`
	RedirectURI         *string `json:"redirect_uri"`
	ClientID            string  `json:"client_id"`
	Scope               string  `json:"scope"`
	State               string  `json:"state"`
	CodeChallenge       string  `json:"code_challenge"`
	CodeChallengeMethod string  `json:"code_challenge_method"`
}

type authCredentials struct {
	Username string `json:"username"`
	Password string `json:"_"`
}

func parseAuthorizeParams(r *http.Request) (*params, error) {
	var redirectURI *string

	if v, ok := r.URL.Query()["redirect_uri"]; ok {
		redirectURI = ptr.String(v[0])
	}

	p := &params{
		ResponseType:        r.URL.Query().Get("response_type"),
		RedirectURI:         redirectURI,
		ClientID:            r.URL.Query().Get("client_id"),
		Scope:               r.URL.Query().Get("scope"),
		State:               r.URL.Query().Get("state"),
		CodeChallenge:       r.URL.Query().Get("code_challenge"),
		CodeChallengeMethod: r.URL.Query().Get("code_challenge_method"),
	}

	if !p.validate() {
		return nil, cher.New("invalid_input", cher.M{
			"params": p,
		})
	}

	return p, nil
}

func parseAuthorizeResponse(r *http.Request) (*authCredentials, *params, error) {
	var redirectURI *string

	if v, ok := r.Form["redirect_uri"]; ok {
		redirectURI = ptr.String(v[0])
	}

	p := &params{
		ResponseType:        r.FormValue("response_type"),
		RedirectURI:         redirectURI,
		ClientID:            r.FormValue("client_id"),
		Scope:               r.FormValue("scope"),
		State:               r.FormValue("state"),
		CodeChallenge:       r.FormValue("code_challenge"),
		CodeChallengeMethod: r.FormValue("code_challenge_method"),
	}

	if !p.validate() {
		return nil, nil, cher.New("invalid_input", cher.M{
			"params": p,
		})
	}

	authCred := &authCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	return authCred, p, nil
}

func (p params) validate() bool {
	if p.ResponseType == "" {
		return false
	}
	if p.RedirectURI != nil && *p.RedirectURI == "" {
		return false
	}
	if p.ClientID == "" {
		return false
	}
	if p.Scope == "" {
		return false
	}
	if p.State == "" {
		return false
	}
	if p.CodeChallenge == "" {
		return false
	}
	if p.CodeChallengeMethod == "" {
		return false
	}
	return true
}
