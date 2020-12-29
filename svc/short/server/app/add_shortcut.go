package app

import (
	"context"

	"dfl/lib/cher"
	"dfl/svc/short"
	"dfl/svc/short/server/db"
)

func (a *App) AddShortcut(ctx context.Context, resource *short.Resource, shortcut string) error {
	err := a.DB.Q.FindShortcutConflicts(ctx, []string{shortcut})
	if err != nil {
		return cher.New("shortcuts_already_taken", cher.M{"shortcut": shortcut}, cher.Coerce(err))
	}

	return a.DB.Q.ChangeShortcut(ctx, db.ArrayAdd, resource.ID, shortcut)
}
