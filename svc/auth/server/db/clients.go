package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/cuvva/ksuid-go"
)

func (qw *QueryableWrapper) GetClientByName(ctx context.Context, name string) (string, error) {
	qb := NewQueryBuilder()
	query, values, err := qb.
		Select("id").
		From("clients").
		Where(sq.Eq{
			"name": name,
		}).
		ToSql()
	if err != nil {
		return "", err
	}

	var clientID string

	row := qw.q.QueryRowContext(ctx, query, values...)

	if err := row.Scan(&clientID); err != nil {
		return "", coerceNotFound(err)
	}

	return clientID, nil
}

func (qw *QueryableWrapper) CreateClient(ctx context.Context, name string) (string, error) {
	qb := NewQueryBuilder()
	query, values, err := qb.
		Insert("clients").
		Columns("id", "name").
		Values(ksuid.Generate("client").String(), name).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return "", err
	}

	var id string

	row := qw.q.QueryRowContext(ctx, query, values...)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
