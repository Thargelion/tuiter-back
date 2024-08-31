package user

import (
	"context"
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
