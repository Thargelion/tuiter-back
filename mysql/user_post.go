package mysql

import (
	"context"
	"tuiter.com/api/userpost"
)

type UserPostRepository struct {
}

func (u UserPostRepository) ListByPage(ctx context.Context, pageId string) ([]*userpost.UserPost, error) {
	//TODO implement me
	panic("implement me")
}
