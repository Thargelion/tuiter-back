package tuit

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageID string) ([]*Post, error)
	Create(ctx context.Context, post *Post) error
	AddLike(ctx context.Context, userID int, tuitID int) error
	RemoveLike(ctx context.Context, userID int, tuitID int) error
}
