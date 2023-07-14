package pguserstorage

import (
	"context"

	"tech-tsarka/internal/storage/user/entity"
)

type UserStorage interface {
	Create(ctx context.Context, arg entity.UserCreateInput) (entity.User, error)
	Get(ctx context.Context, id string) (entity.User, error)
	Update(ctx context.Context, id string, arg entity.UserUpdateInput) error
	Delete(ctx context.Context, id string) error
}
