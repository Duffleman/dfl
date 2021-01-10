package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"dfl/lib/ptr"
	"dfl/svc/short"
)

func (a *App) GetLatestXKCD(ctx context.Context) (*short.Resource, error) {
	comic, err := a.XKCD.Latest(ctx)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006 002", fmt.Sprintf("%d %.3d", comic.Year, comic.Day))
	if err != nil {
		return nil, err
	}

	return &short.Resource{
		ID:        strconv.Itoa(comic.Number),
		Type:      "xkcd",
		Hash:      ptr.String(fmt.Sprintf("XKCD#%d", comic.Day)),
		Name:      &comic.Title,
		Link:      comic.ImageURL,
		MimeType:  ptr.String("image/png"),
		Shortcuts: []string{":xkcd"},
		CreatedAt: date,
	}, nil
}
