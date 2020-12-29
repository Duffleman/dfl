package app

import (
	"context"
	"time"

	authlib "dfl/lib/auth"
	"dfl/lib/cher"
	"dfl/svc/auth"
	"dfl/svc/auth/server/lib/middleware"

	"github.com/cuvva/ksuid-go"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := a.DB.Q.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user.Password == nil || user.InviteCode != nil {
		return nil, cher.New("pending_invite", nil)
	}

	if !checkPasswordHash(req.Password, *user.Password) {
		return nil, cher.New(cher.Unauthorized, nil)
	}

	if !authlib.Can("dflauth:*", user.Scopes) {
		return nil, cher.New(cher.AccessDenied, nil)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, middleware.DFLClaims{
		user.Scopes,
		user.Username,
		jwt.StandardClaims{
			Id:        ksuid.Generate("accesstoken").String(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth.dfl.mn",
			NotBefore: time.Now().Add(-1 * time.Second).Unix(),
			Subject:   user.ID,
			Audience:  "auth.dfl.mn",
		},
	})

	tokenString, err := token.SignedString(a.SK.Private())
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		UserID:    user.ID,
		AuthToken: tokenString,
	}, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
