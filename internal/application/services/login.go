package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/security"
)

func NewUserAuthenticator(userRepo user.LoginRepository, tokenHandler security.TokenHandler) *UserAuthenticator {
	return &UserAuthenticator{userRepo: userRepo, tokenHandler: tokenHandler}
}

func (ua *UserAuthenticator) Login(ctx context.Context, login *user.User) (*user.Logged, error) {
	storedUser, err := ua.userRepo.FindByEmail(ctx, login.Email)

	if err != nil {
		return nil, fmt.Errorf("syserror finding user by email: %w", err)
	}

	err = storedUser.CheckPassword(login.Password)

	if err != nil {
		return nil, fmt.Errorf("wrong password: %w", err)
	}

	token, err := ua.tokenHandler.GenerateToken(storedUser.ID, storedUser.Email, storedUser.Email)

	if err != nil {
		return nil, fmt.Errorf("syserror generating token: %w", err)
	}

	return &user.Logged{User: *storedUser, Token: token}, nil
}

type UserAuthenticator struct {
	userRepo     user.LoginRepository
	tokenHandler security.TokenHandler
}
