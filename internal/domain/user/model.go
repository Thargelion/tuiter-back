package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"tuiter.com/api/pkg/syserror"
)

type Logged struct {
	User
	Token string `json:"token"`
}

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Password  string `json:"password"`
}

// SecureUser will generate an encrypted password. Will remain empty if
// user has no password at all
func (u *User) SecureUser() (*User, error) {
	var pass string
	if u.Password != "" {
		binaryPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("syserror hashing password: %w", err)
		}
		pass = string(binaryPassword)
	} else {
		pass = u.Password
	}

	return &User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Password:  pass,
	}, nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("%w: %w", syserror.ErrUnauthorized, err)
	}

	return nil
}
