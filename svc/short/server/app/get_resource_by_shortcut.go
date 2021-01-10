package app

import (
	"context"

	"dfl/svc/short"
)

// GetResourceByShortcut returns a resource from a :shortcut
func (a *App) GetResourceByShortcut(ctx context.Context, shortcut string) (*short.Resource, error) {
	if shortcut == "xkcd" {
		return a.GetLatestXKCD(ctx)
	}

	return a.DB.Q.FindResourceByShortcut(ctx, shortcut)
}
