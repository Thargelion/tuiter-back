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
	var txResult kit.DatabaseActions
	var res []*user.User
	if len(query) == 0 {
		txResult = r.database.Find(&res)
	} else {
		txResult = r.database.Search(&res, query)
	}
	return res, txResult.Error()
}

func (r *UserRepository) FindUserByID(ctx context.Context, ID string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "id = ?", ID)
	return res, txResult.Error()
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	res := r.database.Create(user)
	return user, res.Error()
}

func NewUserRepository(creator kit.DatabaseActions) *UserRepository {
	return &UserRepository{database: creator}
}
