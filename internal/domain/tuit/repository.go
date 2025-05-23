package tuit

import (
	"context"
)

type LikesRepository interface {
	AddLike(ctx context.Context, userID uint, tuitID int) error
	RemoveLike(ctx context.Context, userID uint, tuitID int) error
}

type Repository interface {
	ListByPage(ctx context.Context, pageID string) ([]*Tuit, error)
	Create(ctx context.Context, post *Tuit) error
	LikesRepository
}

type ReplyRepository interface {
	ReplyListByPage(_ context.Context, parentID uint, pageID int) ([]*Tuit, error)
	Create(ctx context.Context, post *Tuit) error
	LikesRepository
}
