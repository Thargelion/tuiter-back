package tuitfeed

import (
	"context"

	"tuiter.com/api/pkg/query"
)

type Liker interface {
	AddLike
	RemoveLike
}

type AddLike interface {
	AddLike(ctx context.Context, userID uint, tuitID int) (*Model, error)
}

type RemoveLike interface {
	RemoveLike(ctx context.Context, userID uint, tuitID int) (*Model, error)
}

type UseCases interface {
	Paginate(ctx context.Context, userID uint, page int, params query.Params) ([]*Model, error)
	PaginateReplies(ctx context.Context, userID uint, tuitID uint, page int) ([]*Model, error)
	GetByID(ctx context.Context, userID uint, tuitID int) (*Model, error)
}
