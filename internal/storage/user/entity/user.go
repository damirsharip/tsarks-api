package entity

import "errors"

var (
	ErrNotFound = errors.New("product not found")
)

type UserCreateInput struct {
	FirstName string
	LastName  string
}

type UserUpdateInput struct {
	FirstName *string
	LastName  *string
}

type User struct {
	ID        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
