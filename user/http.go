package user

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"net/http"
	"tuiter.com/api/rest"

	"github.com/go-chi/render"
)

type Router struct {
	repo Repository
}

func (r *Router) Search(writer http.ResponseWriter, request *http.Request) {
	var user Filter
	var decoder = schema.NewDecoder()
	var query map[string]interface{}
	queryValues := request.URL.Query()
	err := decoder.Decode(&user, queryValues)
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	data, _ := json.Marshal(user)
	_ = json.Unmarshal(data, &query)
	users, err := r.repo.Search(request.Context(), query)
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
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

func (r *Router) FindUserByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	user, err := r.repo.FindUserByKey(request.Context(), "id", id)
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

func newUserList(users []*User) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, user := range users {
		list = append(list, &Payload{user})
	}

	return list
}

type Filter struct {
	Name *string `json:"name,omitempty"`
}
