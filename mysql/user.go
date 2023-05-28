package mysql

import (
	"context"
	"tuiter.com/api/kit"
	"tuiter.com/api/user"
)

type UserRepository struct {
	database kit.DatabaseActions
}

func (r *UserRepository) Search(ctx context.Context, query map[string]interface{}) ([]*user.User, error) {
	var res []*user.User
	txResult := r.database.Search(&res, query)
	return res, txResult.Error()
}

func (r *UserRepository) FindUserByKey(ctx context.Context, key string, value string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "? = ?", key, value)
	return res, txResult.Error()
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	res := r.database.Create(user)
	return res.Error()
}

func NewUserRepository(creator kit.DatabaseActions) *UserRepository {
	return &UserRepository{database: creator}
}
