package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type DFLClaims struct {
	Version  string   `json:"v"`
	Scopes   []string `json:"scopes"`
	Username string   `json:"username"`
	jwt.StandardClaims
}
