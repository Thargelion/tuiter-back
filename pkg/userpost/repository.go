package userpost

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageNumber int, userID int) ([]*UserPost, error)
}
