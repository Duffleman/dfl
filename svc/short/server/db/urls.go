package db

import (
	"context"
	"fmt"
	"strings"

	"dfl/svc/short"
)

// NewURL inserts a new URL to the database
func (qw *QueryableWrapper) NewURL(ctx context.Context, id, owner, url string) (*short.Resource, error) {
	b := NewQueryBuilder()

	query, values, err := b.
		Insert("resources").
		Columns("id, type, owner, link").
		Values(id, "url", owner, url).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(resourceColumns, ","))).
		ToSql()
	if err != nil {
		return nil, err
	}

	return qw.queryOne(ctx, query, values)
}
