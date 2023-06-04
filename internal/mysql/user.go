package mysql

import (
	"context"
	"fmt"

	"tuiter.com/api/pkg/user"
)

func NewUserRepository(creator DatabaseActions) *UserRepository {
	return &UserRepository{database: creator}
}

type UserRepository struct {
	database DatabaseActions
}

func (r *UserRepository) Search(_ context.Context, query map[string]interface{}) ([]*user.User, error) {
	var txResult DatabaseActions

	var res []*user.User

	if len(query) == 0 {
		txResult = r.database.Find(&res)
	} else {
		txResult = r.database.Search(&res, query)
	}

	if txResult.Error() != nil {
		return nil, fmt.Errorf("error searching users on database %w", txResult.Error())
	}

	return res, nil
}

func (r *UserRepository) FindUserByID(_ context.Context, iD string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "id = ?", iD)

	return res, fmt.Errorf("error finding user on database %s %w", iD, txResult.Error())
}

func (r *UserRepository) Create(_ context.Context, user *user.User) (*user.User, error) {
	txResult := r.database.Create(user)

	if txResult.Error() != nil {
		return nil, fmt.Errorf("error creating user on database %w", txResult.Error())
	}

	return user, nil
}
