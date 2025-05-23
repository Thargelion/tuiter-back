package tuitpost

import (
	"context"

	"tuiter.com/api/pkg/query"
)

type Repository interface {
	SearchByPage(ctx context.Context, userID uint, page int, params query.Params) ([]*TuitPost, error)
	RepliesByPage(ctx context.Context, userID uint, tuitID uint, page int) ([]*TuitPost, error)
	GetByID(ctx context.Context, userID uint, postID int) (*TuitPost, error)
}
