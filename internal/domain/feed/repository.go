package feed

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, userID uint, page int) ([]*Feed, error)
	GetByID(ctx context.Context, userID uint, postID int) (*Feed, error)
}
