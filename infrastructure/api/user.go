package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/internal/user"
)

var errInvalidRequest = errors.New("missing required fields")

func NewUserRouter(useCases user.UseCases, errRenderer ErrorRenderer, logger logging.ContextualLogger) *UserRouter {
	return &UserRouter{
		useCases:      useCases,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type UserRouter struct {
	useCases      user.UseCases
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

func (r *UserRouter) Search(writer http.ResponseWriter, request *http.Request) {
	var filter userFilter

	var decoder = schema.NewDecoder()

	var query map[string]interface{}

	queryValues := request.URL.Query()
	err := decoder.Decode(&filter, queryValues)

	if err != nil {
		err := render.Render(writer, request, r.errorRenderer.RenderError(err))
		if err != nil {
			r.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
			return
		}

		return
	}

	rawFilter, _ := json.Marshal(filter) //nolint:errchkjson
	_ = json.Unmarshal(rawFilter, &query)
	users, err := r.useCases.Search(request.Context(), query)

	if err != nil {
		err := render.Render(writer, request, r.errorRenderer.RenderError(err))
		if err != nil {
			r.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
			return
		}

		return
	}

	err = render.RenderList(writer, request, newUserList(users))
	if err != nil {
		r.logger.Printf(request.Context(), "syserror rendering user list: %v", err)
		return
	}
}

func (r *UserRouter) FindUserByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userFound, err := r.useCases.FindUserByID(request.Context(), id)

	if err != nil {
		err := render.Render(writer, request, r.errorRenderer.RenderError(err))
		if err != nil {
			r.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
			return
		}

		return
	}

	err = render.Render(writer, request, newUserPayload(userFound))
	if err != nil {
		r.logger.Printf(request.Context(), "syserror rendering user: %v", err)
		return
	}
}

func (r *UserRouter) CreateUser(writer http.ResponseWriter, request *http.Request) {
	payload := &userCreatePayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, r.errorRenderer.RenderError(err))
		if err != nil {
			return
		}

		return
	}

	newUser, err := r.useCases.Create(request.Context(), payload.ToUser())
	if err != nil {
		err := render.Render(writer, request, r.errorRenderer.RenderError(err))
		if err != nil {
			return
		}

		return
	}

	err = render.Render(writer, request, newUserPayload(newUser))
	if err != nil {
		return
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
	commonPayload
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func newUserPayload(user *user.User) *userPayload {
	return &userPayload{
		commonPayload: commonPayload{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}
}

func (u *userCreatePayload) Bind(_ *http.Request) error {
	v := validator.New()
	err := v.Struct(u)

	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

type userFilter struct {
	Name *string `json:"name,omitempty"`
}

func (u *userPayload) Bind(_ *http.Request) error {
	return nil
}

func (u *userPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newUserList(users []*user.User) []render.Renderer {
	var list []render.Renderer

	for _, u := range users {
		list = append(list, newUserPayload(u))
	}

	return list
}
