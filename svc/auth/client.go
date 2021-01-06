package auth

import (
	"context"
	"net/http"
	"time"

	"dfl/lib/crpc"
	"dfl/lib/jsonclient"
)

type client struct {
	*crpc.Client
}

func NewClient(baseURL, key string) Service {
	httpClient := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, key),
		Timeout:   5 * time.Second,
	}

	return &client{
		crpc.NewClient(baseURL+"/", httpClient),
	}
}

func (c *client) CreateClient(ctx context.Context, req *CreateClientRequest) (res *CreateClientResponse, err error) {
	return res, c.Do(ctx, "create_client", req, &res)
}

func (c *client) CreateInviteCode(ctx context.Context, req *CreateInviteCodeRequest) (res *CreateInviteCodeResponse, err error) {
	return res, c.Do(ctx, "create_invite_code", req, &res)
}

func (c *client) Token(ctx context.Context, req *TokenRequest) (res *TokenResponse, err error) {
	return res, c.Do(ctx, "token", req, &res)
}

func (c *client) WhoAmI(ctx context.Context) (res *WhoAmIResponse, err error) {
	return res, c.Do(ctx, "whoami", nil, &res)
}
