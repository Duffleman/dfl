package db

import (
	"context"
	"time"

	"dfl/svc/auth"

	sq "github.com/Masterminds/squirrel"
)

func (qw *QueryableWrapper) findOneUser(ctx context.Context, field, value string) (*auth.User, error) {
	qb := NewQueryBuilder()
	query, values, err := qb.
		Select("u.id, u.username, u.email, u.created_at, u.invite_code, u.invite_expiry, u.invite_redeemed_at, u.scopes").
		From("users u").
		Where(sq.Eq{
			field: value,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u auth.User

	row := qw.q.QueryRowContext(ctx, query, values...)

	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.InviteCode, &u.InviteExpiry, &u.InviteRedeemedAt, &u.Scopes); err != nil {
		return nil, coerceNotFound(err)
	}

	return &u, nil
}

func (qw *QueryableWrapper) FindUser(ctx context.Context, id string) (*auth.User, error) {
	return qw.findOneUser(ctx, "id", id)
}

func (qw *QueryableWrapper) GetUserByName(ctx context.Context, username string) (*auth.User, error) {
	return qw.findOneUser(ctx, "username", username)
}

func (qw *QueryableWrapper) RedeemInvite(ctx context.Context, userID string) error {
	qb := NewQueryBuilder()
	query, values, err := qb.
		Update("users u").
		Set("invite_expiry", nil).
		Set("invite_redeemed_at", time.Now()).
		Where(sq.Eq{
			"id": userID,
		}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = qw.q.ExecContext(ctx, query, values...); err != nil {
		return err
	}

	return nil
}
