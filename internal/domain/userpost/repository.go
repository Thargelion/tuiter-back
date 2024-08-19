package userpost

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, userID int, page int) ([]*UserPost, error)
	GetByID(ctx context.Context, userID int, postID int) (*UserPost, error)
}
