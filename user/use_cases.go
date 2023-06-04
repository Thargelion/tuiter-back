package user

import (
	"context"
	"tuiter.com/api/avatar"
)

type UseCases interface {
	CreateUserUseCase
}

type CreateUserUseCase interface {
	Create(ctx context.Context, user *User) (*User, error)
}

type userService struct {
	userRepo       Repository
	generateAvatar avatar.AddAvatarUseCase
}

func (c *userService) Create(ctx context.Context, user *User) (*User, error) {
	user.AvatarURL = c.generateAvatar.New(user.Name)
	return c.userRepo.Create(ctx, user)
}
func NewUserUseCases(userRepo Repository, avatarService avatar.UseCases) UseCases {
	return &userService{userRepo: userRepo, generateAvatar: avatarService}
}
