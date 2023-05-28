package post

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Post, error)
	Create(ctx context.Context, post *Post) error
}