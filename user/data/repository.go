package data

import (
	"context"
	"tuiter.com/api/user/domain"
)

type Repository interface {
	FindUserByID(ctx context.Context, ID string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Search(ctx context.Context, query map[string]interface{}) ([]*domain.User, error)
}
