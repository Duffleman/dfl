package app

import (
	"context"
	"crypto/sha256"
	"time"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	dfljwt "dfl/lib/jwt"
	"dfl/svc/auth"
	"dfl/svc/auth/server/db"

	"github.com/cuvva/ksuid-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/dvsekhvalnov/jose2go/base64url"
)

const atExpiry = 365 * 24 * time.Hour // 365 days

func (a *App) Token(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	client, err := a.DB.Q.FindClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	ac, err := a.DB.Q.FindAuthorizationCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	if time.Now().After(ac.ExpiresAt) {
		return nil, cher.New("code_expired", nil)
	}

	user, err := a.DB.Q.FindUser(ctx, ac.UserID)
	if err != nil {
		return nil, err
	}

	if !authlib.Can(ac.Scope, user.Scopes) {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	switch ac.CodeChallengeMethod {
	case "S256":
		h := sha256.New()
		h.Write([]byte(req.CodeVerifier))
		compare := base64url.Encode(h.Sum(nil))

		if compare != ac.CodeChallenge {
			return nil, cher.New("code_challenge_failed", nil)
		}
	default:
		return nil, cher.New("unsupported_challenge_method", nil)
	}

	expiresAt := time.Now().Add(atExpiry)
	atID := ksuid.Generate("accesstoken").String()

	token := jwt.NewWithClaims(jwt.SigningMethodES384, dfljwt.DFLClaims{
		Scope:    user.Scopes,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			Id:        atID,
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth.dfl.mn",
			NotBefore: time.Now().Add(-1 * time.Second).Unix(),
			Subject:   user.ID,
			Audience:  client.ID,
		},
	})

	tokenString, err := token.SignedString(a.SK.Private())
	if err != nil {
		return nil, err
	}

	err = a.DB.Q.CreateAccessToken(ctx, &db.AccessToken{
		ID:                atID,
		Token:             tokenString,
		AuthorizationCode: ac.ID,
		ExpiresAt:         expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &auth.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		Expires:     int(expiresAt.Sub(time.Now()).Seconds()),
	}, nil
}
