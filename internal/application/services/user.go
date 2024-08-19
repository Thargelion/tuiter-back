package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/avatar"
	"tuiter.com/api/internal/domain/user"
)

type Service struct {
	userRepo       user.Repository
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
	u, err := c.userRepo.FindUserByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("syserror searching for a user on repository: %w", err)
	}

	return u, nil
}

func (c *Service) Create(ctx context.Context, user *user.User) (*user.User, error) {
	user.AvatarURL = c.generateAvatar.New(user.Name)

	newUser, err := c.userRepo.Create(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("syserror creating a user on repository: %w", err)
	}

	return newUser, nil
}

func NewUserUseCases(userRepo user.Repository, avatarService avatar.UseCases) *Service {
	return &Service{userRepo: userRepo, generateAvatar: avatarService}
}
