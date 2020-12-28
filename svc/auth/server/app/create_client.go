package app

import (
	"context"

	"dfl/lib/cher"
	"dfl/svc/auth"
)

func (a *App) CreateClient(ctx context.Context, req *auth.CreateClientRequest) (*auth.CreateClientResponse, error) {
	clientID, err := a.DB.Q.GetClientByName(ctx, req.Name)
	if err != nil {
		if v, ok := err.(cher.E); ok {
			if v.Code == cher.NotFound {
				clientID, err = a.DB.Q.CreateClient(ctx, req.Name)
				if err != nil {
					return nil, err
				}
			}
		}

		return nil, err
	}

	return &auth.CreateClientResponse{
		ClientID: clientID,
	}, nil
}
