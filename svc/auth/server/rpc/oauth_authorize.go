package rpc

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/lib/ptr"
	"dfl/lib/rpc"
	"dfl/lib/slicecontains"
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

		client, err := a.FindClient(r.Context(), params.ClientID)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if err = a.AuthCodeNoNonceExists(r.Context(), params.Nonce); err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if params.RedirectURI != nil && !slicecontains.String(client.RedirectURIs, *params.RedirectURI) {
			rpc.HandleError(w, r, cher.New("invalid_redirect_uri", nil))
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

		client, err := a.FindClient(r.Context(), params.ClientID)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if err = a.AuthCodeNoNonceExists(r.Context(), params.Nonce); err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if params.RedirectURI != nil && !slicecontains.String(client.RedirectURIs, *params.RedirectURI) {
			rpc.HandleError(w, r, cher.New("invalid_redirect_uri", nil))
			return
		}

		user, err := a.GetUserByName(r.Context(), authCredentials.Username)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if !authlib.Can(params.Scope, user.Scopes) {
			rpc.HandleError(w, r, cher.New(cher.AccessDenied, nil, cher.New("invalid_scopes", nil)))
			return
		}

		res, err := a.Authorization(r.Context(), &auth.AuthorizationRequest{
			ClientID:            client.ID,
			ResponseType:        params.ResponseType,
			RedirectURI:         params.RedirectURI,
			State:               params.State,
			Nonce:               params.Nonce,
			CodeChallengeMethod: params.CodeChallengeMethod,
			CodeChallenge:       params.CodeChallenge,
			Username:            authCredentials.Username,
			Password:            authCredentials.Password,
			Scope:               params.Scope,
		}, user)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		if params.RedirectURI == nil {
			displayAuthToken(w, r, res, client)
			return
		}

		urlVals := &url.Values{
			"code":  []string{res.AuthorizationCode},
			"state": []string{res.State},
		}

		url := fmt.Sprintf("%s?%s", *params.RedirectURI, urlVals.Encode())

		w.Header().Set("Content-Type", "")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func displayAuthToken(w http.ResponseWriter, r *http.Request, res *auth.AuthorizationResponse, client *auth.Client) {
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

type params struct {
	ResponseType        string  `json:"response_type"`
	RedirectURI         *string `json:"redirect_uri"`
	ClientID            string  `json:"client_id"`
	Scope               string  `json:"scope"`
	State               string  `json:"state"`
	Nonce               string  `json:"nonce"`
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
		Nonce:               r.URL.Query().Get("nonce"),
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

	if v := r.FormValue("redirect_uri"); v != "" {
		redirectURI = ptr.String(v)
	}

	p := &params{
		ResponseType:        r.FormValue("response_type"),
		RedirectURI:         redirectURI,
		ClientID:            r.FormValue("client_id"),
		Scope:               r.FormValue("scope"),
		State:               r.FormValue("state"),
		Nonce:               r.FormValue("nonce"),
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
	if !slicecontains.String([]string{"code"}, p.ResponseType) {
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
	if p.Nonce == "" {
		return false
	}
	if p.CodeChallenge == "" {
		return false
	}
	if p.CodeChallengeMethod == "" {
		return false
	}
	if !slicecontains.String([]string{"S256"}, p.CodeChallengeMethod) {
		return false
	}
	return true
}
