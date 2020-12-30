package auth

import (
	"context"
	"time"
)

type Service interface {
	CreateClient(context.Context, *CreateClientRequest) (*CreateClientResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Register(context.Context, *RegisterRequest) error
	Token(context.Context, *TokenRequest) (*TokenResponse, error)
	WhoAmI(context.Context, *WhoAmIRequest) (*WhoAmIResponse, error)
}

type User struct {
	ID               string     `json:"id"`
	Username         string     `json:"username"`
	Email            *string    `json:"email"`
	Password         *string    `json:"password"`
	InviteCode       *string    `json:"invite_code"`
	CreatedAt        time.Time  `json:"created_at"`
	InviteExpiry     *time.Time `json:"invite_expiry"`
	InviteRedeemedAt *time.Time `json:"invite_redeemed_at"`
	Scopes           string     `json:"scopes"`
}

type Client struct {
	ID           string    `json:"client_id"`
	Name         string    `json:"name"`
	RedirectURIs []string  `json:"redirect_uris"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
}

type RegisterRequest struct {
	Username   string  `json:"username"`
	Email      *string `json:"email"`
	Password   string  `json:"password"`
	InviteCode string  `json:"invite_code"`
}

type WhoAmIRequest struct {
	Username string `json:"username"`
}

type WhoAmIResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type CreateClientRequest struct {
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris"`
}

type CreateClientResponse struct {
	ClientID string `json:"client_id"`
}

type AuthorizationCode struct {
	ID                  string    `json:"id"`
	ClientID            string    `json:"client_id"`
	ResponseType        string    `json:"response_type"`
	RedirectURI         *string   `json:"redirect_uri"`
	State               string    `json:"state"`
	CodeChallengeMethod string    `json:"code_challenge_method"`
	CodeChallenge       string    `json:"code_challenge"`
	UserID              string    `json:"user_id"`
	CreatedAt           time.Time `json:"created_at"`
	ExpiresAt           time.Time `json:"expires_at"`
	Scope               string    `json:"scope"`
}

type AuthorizationRequest struct {
	ClientID            string  `json:"client_id"`
	ResponseType        string  `json:"response_type"`
	RedirectURI         *string `json:"redirect_uri"`
	State               string  `json:"state"`
	CodeChallengeMethod string  `json:"code_challenge_method"`
	CodeChallenge       string  `json:"code_challenge"`
	Username            string  `json:"username"`
	Password            string  `json:"password"`
	Scope               string  `json:"scope"`
}

type AuthorizationResponse struct {
	AuthorizationCode string `json:"authorization_code"`
	ExpiresAt         string `json:"expires_at"`
	ExpiresIn         int    `json:"expires_in"`
	State             string `json:"state"`
}

type TokenRequest struct {
	ClientID     string  `json:"client_id"`
	GrantType    string  `json:"grant_type"`
	RedirectURI  *string `json:"redirect_uri"`
	Code         string  `json:"code"`
	CodeVerifier string  `json:"code_verifier"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires"`
}
