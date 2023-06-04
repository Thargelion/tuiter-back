package userpost

import (
	"context"
)

type Repository interface {
	ListByPage(ctx context.Context, pageID string) ([]*UserPost, error)
}
