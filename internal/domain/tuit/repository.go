package tuit

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageID string) ([]*Tuit, error)
	Create(ctx context.Context, post *Tuit) error
	AddLike(ctx context.Context, userID uint, tuitID int) error
	RemoveLike(ctx context.Context, userID uint, tuitID int) error
}
