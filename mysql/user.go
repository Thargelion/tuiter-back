package mysql

import (
	"context"
	"gorm.io/gorm"
	"tuiter.com/api/user"
)

type UserRepository struct {
	database *gorm.DB
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "name = ?", username)
	return res, txResult.Error
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	res := r.database.Create(user)
	return res.Error
}

func NewUserRepository(creator *gorm.DB) *UserRepository {
	return &UserRepository{database: creator}
}
