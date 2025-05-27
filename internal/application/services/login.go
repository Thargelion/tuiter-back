package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/security"
	"tuiter.com/api/pkg/syserror"
)

func NewUserAuthenticator(userRepo user.LoginRepository, tokenHandler security.TokenHandler) *UserAuthenticator {
	return &UserAuthenticator{userRepo: userRepo, tokenHandler: tokenHandler}
}

func (ua *UserAuthenticator) Login(ctx context.Context, login *user.User) (*user.Logged, error) {
	storedUser, err := ua.userRepo.FindByEmail(ctx, login.Email)

	if err != nil {
		return nil, fmt.Errorf("%w: wrong user or password", syserror.ErrUnauthorized)
	}

	err = storedUser.CheckPassword(login.Password)

	if err != nil {
		return nil, fmt.Errorf("%w: wrong user or password", syserror.ErrUnauthorized)
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
