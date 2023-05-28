package user

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"tuiter.com/api/rest"

	"github.com/go-chi/render"
)

type Router struct {
	repo Repository
}

func (r *Router) FindUser(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "userName")
	user, err := r.repo.FindUserByUsername(request.Context(), id)
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, &Payload{user})
	if err != nil {
		return
	}
}

func (r *Router) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data := &Payload{}
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

type Payload struct {
	*User
}

func (u *Payload) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("missing required User fields")
	}

	return nil
}

func (u *Payload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
