package user

import (
	"context"
)

type UseCases interface {
	SearchUserUseCase
	FindUserUseCase
	EditUserUseCase
	CreateAndLogin(ctx context.Context, user *User) (*Logged, error)
}

type EditUserUseCase interface {
	Update(ctx context.Context, user *User) (*User, error)
}

type SearchUserUseCase interface {
	Search(ctx context.Context, query map[string]interface{}) ([]*User, error)
}

type FindUserUseCase interface {
	FindUserByID(ctx context.Context, ID string) (*User, error)
}

type Authenticate interface {
	// Login will return a logged user if the login is successful
	Login(ctx context.Context, login *User) (*Logged, error)
}
