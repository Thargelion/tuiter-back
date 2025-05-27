package services

import (
	"context"

	"tuiter.com/api/internal/domain/tuit"
)

func (u *UseReply) ListByPage(ctx context.Context, pageID string) ([]*tuit.Tuit, error) {
	// TODO implement me
	panic("implement me")
}

func (u *UseReply) Create(ctx context.Context, post *tuit.Tuit) error {
	// TODO implement me
	panic("implement me")
}

func (u *UseReply) AddLike(ctx context.Context, userID uint, tuitID int) error {
	// TODO implement me
	panic("implement me")
}

func (u *UseReply) RemoveLike(ctx context.Context, userID uint, tuitID int) error {
	// TODO implement me
	panic("implement me")
}

type UseReply struct {
	tuitRepository tuit.Repository
}
