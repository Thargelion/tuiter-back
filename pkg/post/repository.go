package post

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageID string) ([]*Post, error)
	Create(ctx context.Context, post *Post) error
	AddLike(ctx context.Context, postID int, userID int) error
}
