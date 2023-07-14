package pguserstorage

import (
	"context"

	"tech-tsarka/internal/storage/user/entity"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) Update(ctx context.Context, id string, arg entity.UserUpdateInput) error {
	query := queryBuilder.Update(users)

	if arg.FirstName != nil {
		query = query.Set("first_name", *arg.FirstName)
	}

	if arg.LastName != nil {
		query = query.Set("last_name", *arg.LastName)
	}

	query = query.Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	usr, err := s.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if usr.RowsAffected() == 0 {
		return entity.ErrNotFound
	}

	return nil
}
