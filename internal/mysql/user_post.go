package mysql

import (
	"context"

	"tuiter.com/api/pkg/userpost"
)

type UserPostRepository struct {
}

func (u UserPostRepository) ListByPage(_ context.Context, _ string) ([]*userpost.UserPost, error) {
	panic("implement me")
}
