package data

import (
	"context"
	"tuiter.com/api/post/domain"
)

type Repository interface {
	FindAll(ctx context.Context, pageId string) ([]*domain.Post, error)
	Create(ctx context.Context, post *domain.Post) error
}
