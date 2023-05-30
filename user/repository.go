package user

import (
	"context"
)

type Repository interface {
	FindUserByID(ctx context.Context, ID string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Search(ctx context.Context, query map[string]interface{}) ([]*User, error)
}
