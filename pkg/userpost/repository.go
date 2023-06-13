package userpost

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageNumber int, userID int) ([]*UserPost, error)
	GetByID(ctx context.Context, userID int, postID int) (*UserPost, error)
}
