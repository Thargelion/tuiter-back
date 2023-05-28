package mysql

import (
	"context"
	"tuiter.com/api/kit"
	"tuiter.com/api/user"
)

type UserRepository struct {
	database kit.Dao
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "name = ?", username)
	return res, txResult.Error()
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	res := r.database.Create(user)
	return res.Error()
}

func NewUserRepository(creator kit.Dao) *UserRepository {
	return &UserRepository{database: creator}
}
