package feed

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, userID int, page int) ([]*Feed, error)
	GetByID(ctx context.Context, userID int, postID int) (*Feed, error)
}
