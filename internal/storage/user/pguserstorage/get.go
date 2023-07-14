package pguserstorage

import (
	"context"

	"tech-tsarka/internal/storage/user/entity"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) Get(ctx context.Context, id string) (entity.User, error) {
	query := queryBuilder.
		Select("first_name").
		Column("last_name").
		From(users).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()

	if err != nil {
		return entity.User{}, err
	}

	var usr entity.User

	row := s.db.QueryRow(ctx, sql, args...)

	if err = row.Scan(
		&usr.FirstName,
		&usr.LastName,
	); err != nil {
		return entity.User{}, err
	}

	return usr, nil
}
