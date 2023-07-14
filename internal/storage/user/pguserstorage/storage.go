package pguserstorage

import (
	"context"
	"fmt"

	"tech-tsarka/internal/service/userservice"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var _ userservice.UserStorage = (*storage)(nil)

const (
	scheme     = "public"
	usersTable = "users"
)

var (
	users        = fmt.Sprintf("%s.%s", scheme, usersTable)
	queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type DB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type storage struct {
	db DB
}

func NewStorage(p DB) *storage {
	return &storage{
		db: p,
	}
}
