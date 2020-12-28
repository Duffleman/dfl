package auth

import "time"

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
	Name string `json:"name"`
}

type CreateClientResponse struct {
	ClientID string `json:"client_id"`
}
