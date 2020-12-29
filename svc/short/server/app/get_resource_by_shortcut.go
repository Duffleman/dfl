package app

import (
	"context"

	"dfl/svc/short"
)

// GetResourceByShortcut returns a resource from a :shortcut
func (a *App) GetResourceByShortcut(ctx context.Context, shortcut string) (*short.Resource, error) {
	return a.DB.Q.FindResourceByShortcut(ctx, shortcut)
}
