package api

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/user/domain"
)

type userPayload struct {
	*domain.User
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

func newUserList(users []*domain.User) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, user := range users {
		list = append(list, &userPayload{user})
	}

	return list
}
