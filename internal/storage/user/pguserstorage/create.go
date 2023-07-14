package pguserstorage

import (
	"context"

	"tech-tsarka/internal/storage/user/entity"
)

func (s *storage) Create(ctx context.Context, arg entity.UserCreateInput) (entity.User, error) {
	query := queryBuilder.
		Insert(users).
		Columns("first_name", "last_name").
		Values(arg.FirstName, arg.LastName).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()

	if err != nil {
		return entity.User{}, err
	}

	var id string
	row := s.db.QueryRow(ctx, sql, args...)

	if err = row.Scan(
		&id,
	); err != nil {
		return entity.User{}, err
	}

	return entity.User{
		ID: id,
	}, nil
}
