package tuitfeed

import (
	"context"

	"tuiter.com/api/pkg/query"
)

type Repository interface {
	SearchByPage(ctx context.Context, userID uint, page int, params query.Params) ([]*Model, error)
	RepliesByPage(ctx context.Context, userID uint, tuitID uint, page int) ([]*Model, error)
	GetByID(ctx context.Context, userID uint, postID int) (*Model, error)
}
