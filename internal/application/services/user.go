package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/avatar"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/security"
)

type Service struct {
	userRepo       user.Repository
	tokenHandler   security.TokenHandler
	generateAvatar avatar.AddAvatarUseCase
}

func (c *Service) Search(ctx context.Context, query map[string]interface{}) ([]*user.User, error) {
	users, err := c.userRepo.Search(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("syserror searching for a user on repository: %w", err)
	}

	return users, nil
}

func (c *Service) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	userByID, err := c.userRepo.FindUserByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("syserror searching for a user on repository: %w", err)
	}

	return userByID, nil
}

func (c *Service) create(ctx context.Context, user *user.User) (*user.User, error) {
	user.AvatarURL = c.generateAvatar.New(user.Name)

	newUser, err := c.userRepo.Create(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("syserror creating a user on repository: %w", err)
	}

	return newUser, nil
}

func (c *Service) CreateAndLogin(ctx context.Context, u *user.User) (*user.Logged, error) {
	secureUser, err := u.SecureUser()

	if err != nil {
		return nil, fmt.Errorf("syserror securing user: %w", err)
	}

	newUser, err := c.create(ctx, secureUser)

	if err != nil {
		return nil, fmt.Errorf("syserror creating user: %w", err)
	}

	token, err := c.tokenHandler.GenerateToken(newUser.Email, newUser.Email)

	if err != nil {
		return nil, fmt.Errorf("syserror generating token: %w", err)
	}

	return &user.Logged{User: *newUser, Token: token}, nil
}

func NewUserService(
	userRepo user.Repository,
	tokenHandler security.TokenHandler,
	avatarService avatar.UseCases,
) *Service {
	return &Service{
		userRepo:       userRepo,
		tokenHandler:   tokenHandler,
		generateAvatar: avatarService,
	}
}
