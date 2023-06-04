package user

import (
	"context"
	"fmt"

	"tuiter.com/api/avatar"
)

type UseCases interface {
	CreateUserUseCase
	SearchUserUseCase
	FindUserUseCase
}

type CreateUserUseCase interface {
	Create(ctx context.Context, user *User) (*User, error)
}

type SearchUserUseCase interface {
	Search(ctx context.Context, query map[string]interface{}) ([]*User, error)
}

type FindUserUseCase interface {
	FindUserByID(ctx context.Context, ID string) (*User, error)
}

type Service struct {
	userRepo       Repository
	generateAvatar avatar.AddAvatarUseCase
}

func (c *Service) Search(ctx context.Context, query map[string]interface{}) ([]*User, error) {
	users, err := c.userRepo.Search(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("error searching for a user on repository: %w", err)
	}

	return users, nil
}

func (c *Service) FindUserByID(ctx context.Context, id string) (*User, error) {
	user, err := c.userRepo.FindUserByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error searching for a user on repository: %w", err)
	}

	return user, nil
}

func (c *Service) Create(ctx context.Context, user *User) (*User, error) {
	user.AvatarURL = c.generateAvatar.New(user.Name)

	newUser, err := c.userRepo.Create(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("error creating a user on repository: %w", err)
	}

	return newUser, nil
}
func NewUserUseCases(userRepo Repository, avatarService avatar.UseCases) *Service {
	return &Service{userRepo: userRepo, generateAvatar: avatarService}
}
