package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/logging"
)

func (u *UserEntity) TableName() string {
	return "users"
}

func (u *UserEntity) ToModel() user.User {
	return user.User{
		ID:        u.User.ID,
		Name:      u.User.Name,
		Email:     u.Email,
		AvatarURL: u.User.AvatarURL,
	}
}

type UserEntity struct {
	user.User
	Email string `gorm:"index:idx_email,unique"`
	gorm.Model
}

func NewUserRepository(creator *gorm.DB, logger logging.ContextualLogger) *UserRepository {
	return &UserRepository{database: creator, logger: logger}
}

type UserRepository struct {
	database *gorm.DB
	logger   logging.ContextualLogger
}

func (r *UserRepository) Search(ctx context.Context, query map[string]interface{}) ([]*user.User, error) {
	var txResult *gorm.DB

	var res []*user.User

	if len(query) == 0 {
		txResult = r.database.Find(&res)
	} else {
		txResult = r.database.Find(&res, query)
	}

	if txResult.Error != nil {
		r.logger.Printf(ctx, "syserror searching users on database %v", txResult.Error)

		return nil, fmt.Errorf("syserror searching users on database %w", txResult.Error)
	}

	return res, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, userID string) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "id = ?", userID)

	if txResult.Error != nil {
		r.logger.Printf(ctx, "syserror finding user on database %s %v", userID, txResult.Error)

		return nil, fmt.Errorf("syserror finding user on database %s %w", userID, txResult.Error)
	}

	return res, nil
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	txResult := r.database.Create(user)

	if txResult.Error != nil {
		r.logger.Printf(ctx, "syserror creating user on database %v", txResult.Error)

		return nil, fmt.Errorf("syserror creating user on database %w", txResult.Error)
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(
	_ context.Context,
	email string,
) (*user.User, error) {
	var res = &user.User{}
	txResult := r.database.First(&res, "email = ?", email)

	if txResult.Error != nil {
		return nil, fmt.Errorf("syserror finding user on database %s %w", email, txResult.Error)
	}

	return res, nil
}
