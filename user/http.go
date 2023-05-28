package user

import (
	"errors"
	"net/http"
	"tuiter.com/api/rest"

	"github.com/go-chi/render"
)

type Router struct {
	repo Repository
}

func (r *Router) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data := &Request{}
	if err := render.Bind(request, data); err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err := r.repo.Create(request.Context(), data.User)
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, rest.NewResponse(201, "User created"))
	if err != nil {
		return
	}
}

func NewUserRouter(repository Repository) *Router {
	return &Router{repository}
}

type Request struct {
	*User
}

func (u Request) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("missing required User fields")
	}

	return nil
}
