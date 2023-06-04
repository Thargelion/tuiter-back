package api

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/user"
)

import "github.com/go-playground/validator/v10"

type userCreatePayload struct {
	Name      string `json:"name" validate:"required"`
	AvatarURL string `json:"avatar_url"`
}

func (u *userCreatePayload) ToUser() *user.User {
	return &user.User{
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}
}

type userPayload struct {
	*user.User
}

func (u *userCreatePayload) Bind(r *http.Request) error {
	v := validator.New()
	err := v.Struct(u)
	if err != nil {
		return err
	}

	return nil
}

type userFilter struct {
	Name *string `json:"name,omitempty"`
}

func (u *userPayload) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("missing required User fields")
	}

	return nil
}

func (u *userPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newUserList(users []*user.User) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, u := range users {
		list = append(list, &userPayload{u})
	}

	return list
}
