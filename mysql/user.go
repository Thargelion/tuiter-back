package mysql

import (
	"context"
	"tuiter.com/api/pkg"
	"tuiter.com/api/user"
)

type userRepository struct {
	database pkg.DatabaseActions
}

func (r *userRepository) Search(ctx context.Context, query map[string]interface{}) ([]*user.User, error) {
	var txResult pkg.DatabaseActions
	var res []*user.User
	if len(query) == 0 {
		txResult = r.database.Find(&res)
	} else {
		txResult = r.database.Search(&res, query)
	}
	return res, txResult.Error()
}

func (r *userRepository) FindUserByID(ctx context.Context, ID string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "id = ?", ID)
	return res, txResult.Error()
}

func (r *userRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	res := r.database.Create(user)
	return user, res.Error()
}

func NewUserRepository(creator pkg.DatabaseActions) user.Repository {
	return &userRepository{database: creator}
}
