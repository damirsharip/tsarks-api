package userservice

import (
	"context"
	"errors"
	"fmt"

	userentity "tech-tsarka/internal/storage/user/entity"

	"github.com/google/uuid"
)

var ErrValidation = errors.New("invalid data")

type UserStorage interface {
	Create(ctx context.Context, arg userentity.UserCreateInput) (userentity.User, error)
	Get(ctx context.Context, id string) (userentity.User, error)
	Update(ctx context.Context, id string, arg userentity.UserUpdateInput) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	user UserStorage
}

func NewService(order UserStorage) *service {
	return &service{
		user: order,
	}
}

func (s *service) Create(ctx context.Context, arg userentity.UserCreateInput) (userentity.User, error) {
	if arg.FirstName == "" {
		return userentity.User{}, fmt.Errorf("%w: field: [first_name] cannot be empty", ErrValidation)
	}
	if arg.LastName == "" {
		return userentity.User{}, fmt.Errorf("%w: field: [last_name] cannot be empty", ErrValidation)
	}

	return s.user.Create(ctx, arg)
}

func (s *service) Get(ctx context.Context, id string) (userentity.User, error) {
	if id == "" {
		return userentity.User{}, fmt.Errorf("%w: field: [id] cannot be empty", ErrValidation)
	}

	if _, err := uuid.Parse(id); err != nil {
		return userentity.User{}, fmt.Errorf("%w: field: [id] is invalid uuid", ErrValidation)
	}

	return s.user.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, arg userentity.UserUpdateInput) error {
	if id == "" {
		return fmt.Errorf("%w: field: [id] cannot be empty", ErrValidation)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: field: [id] is invalid uuid", ErrValidation)
	}

	if arg.FirstName == nil && arg.LastName == nil {
		return fmt.Errorf("%w: no fields to update", ErrValidation)
	}

	return s.user.Update(ctx, id, arg)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: field: [id] cannot be empty", ErrValidation)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: field: [id] is invalid uuid", ErrValidation)
	}

	return s.user.Delete(ctx, id)
}
