package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"net/http"
	"tuiter.com/api/user"
)

type userRouter struct {
	useCases user.UseCases
	repo     user.Repository
}

func (r *userRouter) Search(writer http.ResponseWriter, request *http.Request) {
	var filter userFilter
	var decoder = schema.NewDecoder()
	var query map[string]interface{}
	queryValues := request.URL.Query()
	err := decoder.Decode(&filter, queryValues)
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	rawFilter, _ := json.Marshal(filter)
	_ = json.Unmarshal(rawFilter, &query)
	users, err := r.repo.Search(request.Context(), query)
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.RenderList(writer, request, newUserList(users))
	if err != nil {
		return
	}
}

func (r *userRouter) FindUserByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userFound, err := r.repo.FindUserByID(request.Context(), id)
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, &userPayload{userFound})
	if err != nil {
		return
	}
}

func (r *userRouter) CreateUser(writer http.ResponseWriter, request *http.Request) {
	payload := &userCreatePayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	newUser, err := r.useCases.Create(request.Context(), payload.ToUser())
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			fmt.Println(err)
			return
		}

		return
	}

	err = render.Render(writer, request, &userPayload{newUser})
	if err != nil {
		return
	}
}

func NewUserRouter(repository user.Repository, useCases user.UseCases) user.Api {
	return &userRouter{
		repo:     repository,
		useCases: useCases,
	}
}

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
