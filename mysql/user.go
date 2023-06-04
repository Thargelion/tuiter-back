package mysql

import (
	"context"
	"tuiter.com/api/kit"
	"tuiter.com/api/user/domain"
)

type UserRepository struct {
	database kit.DatabaseActions
}

func (r *UserRepository) Search(ctx context.Context, query map[string]interface{}) ([]*domain.User, error) {
	var txResult kit.DatabaseActions
	var res []*domain.User
	if len(query) == 0 {
		txResult = r.database.Find(&res)
	} else {
		txResult = r.database.Search(&res, query)
	}
	return res, txResult.Error()
}

func (r *UserRepository) FindUserByID(ctx context.Context, ID string) (*domain.User, error) {
	var res = &domain.User{}
	txResult := r.database.First(&res, "id = ?", ID)
	return res, txResult.Error()
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	res := r.database.Create(user)
	return user, res.Error()
}

func NewUserRepository(creator kit.DatabaseActions) *UserRepository {
	return &UserRepository{database: creator}
}
