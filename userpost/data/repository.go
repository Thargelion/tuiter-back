package data

import (
	"context"
	"tuiter.com/api/userpost"
)

type UserPostRepository interface {
	ListByPage(ctx context.Context, pageId string) ([]*userpost.UserPost, error)
}
