package user

import (
	"context"
)

type Repository interface {
	FindUserByKey(ctx context.Context, key string, value string) (*User, error)
	Create(ctx context.Context, user *User) error
	Search(ctx context.Context, query map[string]interface{}) ([]*User, error)
}
