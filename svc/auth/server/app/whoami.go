package app

import (
	"context"

	"dfl/svc/auth"
)

func (a *App) WhoAmI(ctx context.Context, req *auth.WhoAmIRequest) (*auth.WhoAmIResponse, error) {
	user, err := a.DB.Q.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &auth.WhoAmIResponse{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}
