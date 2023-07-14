package pguserstorage

import (
	"context"

	"tech-tsarka/internal/storage/user/entity"

	sq "github.com/Masterminds/squirrel"
)

func (s *storage) Delete(ctx context.Context, id string) error {
	query := queryBuilder.Delete(users).
		Where(sq.Eq{"id": id})

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
