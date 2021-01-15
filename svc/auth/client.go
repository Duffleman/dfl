package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type client struct {
	*crpc.Client
}

func NewClient(baseURL, key string) Service {
	httpClient := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, "Bearer", key),
		Timeout:   5 * time.Second,
	}

	return &client{
		crpc.NewClient(baseURL+"/", httpClient),
	}
}

func (c *client) CreateClient(ctx context.Context, req *CreateClientRequest) (res *CreateClientResponse, err error) {
	return res, c.Do(ctx, "create_client", "2021-01-15", req, &res)
}

func (c *client) CreateInviteCode(ctx context.Context, req *CreateInviteCodeRequest) (res *CreateInviteCodeResponse, err error) {
	return res, c.Do(ctx, "create_invite_code", "2021-01-15", req, &res)
}

func (c *client) Token(ctx context.Context, req *TokenRequest) (res *TokenResponse, err error) {
	return res, c.Do(ctx, "token", "2021-01-15", req, &res)
}

func (c *client) WhoAmI(ctx context.Context) (res *WhoAmIResponse, err error) {
	return res, c.Do(ctx, "whoami", "2021-01-15", nil, &res)
}
