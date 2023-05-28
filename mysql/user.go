package mysql

import (
	"context"
	"tuiter.com/api/kit"
	"tuiter.com/api/user"
)

type UserRepository struct {
	creator kit.Creator
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	res := r.creator.Create(user)
	return res.Error
}

func NewUserRepository(creator kit.Creator) *UserRepository {
	return &UserRepository{creator: creator}
}
