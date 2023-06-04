package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"net/http"
	"tuiter.com/api/api"
	"tuiter.com/api/user"
)

type UserRouter struct {
	useCases user.UseCases
	repo     user.Repository
}

func (r *UserRouter) Search(writer http.ResponseWriter, request *http.Request) {
	var filter userFilter
	var decoder = schema.NewDecoder()
	var query map[string]interface{}
	queryValues := request.URL.Query()
	err := decoder.Decode(&filter, queryValues)
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	rawFilter, _ := json.Marshal(filter)
	_ = json.Unmarshal(rawFilter, &query)
	users, err := r.repo.Search(request.Context(), query)
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
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

func (r *UserRouter) FindUserByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	user, err := r.repo.FindUserByID(request.Context(), id)
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, &userPayload{user})
	if err != nil {
		return
	}
}

func (r *UserRouter) CreateUser(writer http.ResponseWriter, request *http.Request) {
	payload := &userCreatePayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	newUser, err := r.useCases.Create(request.Context(), payload.ToUser())
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
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

func NewUserRouter(repository user.Repository, useCases user.UseCases) *UserRouter {
	return &UserRouter{
		repo:     repository,
		useCases: useCases,
	}
}
