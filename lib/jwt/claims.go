package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type DFLClaims struct {
	Scope    string `json:"scope"`
	Username string `json:"username"`
	jwt.StandardClaims
}
