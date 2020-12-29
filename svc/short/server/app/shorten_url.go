package app

import (
	"context"
	"fmt"
	"time"

	"dfl/svc/short"

	"github.com/cuvva/ksuid-go"
)

// ShortenURL shortens a URL
func (a *App) ShortenURL(ctx context.Context, url, username string) (*short.CreateResourceResponse, error) {
	urlID := ksuid.Generate("url").String()

	// save to DB
	urlRes, err := a.DB.Q.NewURL(ctx, urlID, username, url)
	if err != nil {
		return nil, err
	}

	hash := a.makeHash(urlRes.Serial)
	fullURL := fmt.Sprintf("%s/%s", a.RootURL, hash)

	gctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go a.saveHash(gctx, cancel, urlRes.Serial, hash)

	return &short.CreateResourceResponse{
		ResourceID: urlRes.ID,
		Type:       urlRes.Type,
		Hash:       hash,
		URL:        fullURL,
	}, nil
}
