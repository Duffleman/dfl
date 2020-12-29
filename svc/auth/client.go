package auth

import (
	"context"
	"net/http"
	"time"

	"dfl/lib/crpc"
)

type client struct {
	*crpc.Client
}

func NewClient(baseURL string) Service {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &client{
		crpc.NewClient(baseURL+"/", httpClient),
	}
}

func (c *client) CreateClient(ctx context.Context, req *CreateClientRequest) (res *CreateClientResponse, err error) {
	return res, c.Do(ctx, "create_client", req, &res)
}

func (c *client) Login(ctx context.Context, req *LoginRequest) (res *LoginResponse, err error) {
	return res, c.Do(ctx, "login", req, &res)
}

func (c *client) Register(ctx context.Context, req *RegisterRequest) error {
	return c.Do(ctx, "register", req, nil)
}

func (c *client) Token(ctx context.Context, req *TokenRequest) (res *TokenResponse, err error) {
	return res, c.Do(ctx, "token", req, &res)
}

func (c *client) WhoAmI(ctx context.Context, req *WhoAmIRequest) (res *WhoAmIResponse, err error) {
	return res, c.Do(ctx, "whoami", req, &res)
}
