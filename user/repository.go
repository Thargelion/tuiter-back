package user

import "context"

type Repository interface {
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}
