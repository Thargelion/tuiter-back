package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

var errInvalidRequest = errors.New("missing required fields")

func NewUserHandler(
	useCases user.UseCases,
	userExtractor security.UserExtractor,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *User {
	return &User{
		useCases:      useCases,
		userExtractor: userExtractor,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type User struct {
	useCases      user.UseCases
	userExtractor security.UserExtractor
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

func (u *User) Search(writer http.ResponseWriter, request *http.Request) {
	var filter userFilter

	var decoder = schema.NewDecoder()

	var query map[string]interface{}

	queryValues := request.URL.Query()
	err := decoder.Decode(&filter, queryValues)

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	rawFilter, _ := json.Marshal(filter) //nolint:errchkjson
	_ = json.Unmarshal(rawFilter, &query)
	users, err := u.useCases.Search(request.Context(), query)

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	err = render.RenderList(writer, request, newUserList(users))
	if err != nil {
		u.logger.Printf(request.Context(), "syserror rendering user list: %v", err)

		return
	}
}

// FindUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} userPayload
// @Router /users/{id} [get].
func (u *User) FindUserByID(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	userFound, err := u.useCases.FindUserByID(request.Context(), id)

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	err = render.Render(writer, request, newUserPayload(userFound))
	if err != nil {
		u.logger.Printf(request.Context(), "syserror rendering user: %v", err)

		return
	}
}

// MeUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} userPayload
// @Router /me [get].
func (u *User) MeUser(writer http.ResponseWriter, request *http.Request) {
	token, ok := request.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(writer, request, ErrInvalidRequest(errors.New("unauthorized")))

		return
	}

	userId, err := u.userExtractor.ExtractUserId(token)
	userFound, err := u.useCases.FindUserByID(request.Context(), strconv.FormatInt(int64(userId), 10))

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	err = render.Render(writer, request, newUserPayload(userFound))
	if err != nil {
		u.logger.Printf(request.Context(), "syserror rendering user: %v", err)

		return
	}
}

// UpdateProfile Update user godoc
// @Summary create a new user
// @Description create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body updateUserPayload true "User"
// @Success 200 {object} loggedUserPayload
// @Router /me/profile [put].
func (u *User) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	payload := &userEditPayload{}
	if err := render.Bind(r, payload); err != nil {
		err := render.Render(w, r, u.errorRenderer.RenderError(err))
		if err != nil {
			return
		}

		return
	}
	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("unauthorized")))

		return
	}

	model := payload.ToUser()

	id, err := u.userExtractor.ExtractUserId(token)

	model.ID = id

	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	updated, err := u.useCases.Update(r.Context(), model)

	if err != nil {
		_ = render.Render(w, r, u.errorRenderer.RenderError(err))
	}

	err = render.Render(w, r, newUserPayload(updated))
	if err != nil {
		u.logger.Printf(r.Context(), "syserror rendering user: %v", err)

		return
	}
}

// CreateUser create Users godoc
// @Summary create a new user
// @Description create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body userCreatePayload true "User"
// @Success 200 {object} loggedUserPayload
// @Router /users [post].
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload := &userCreatePayload{}
	if err := render.Bind(r, payload); err != nil {
		err := render.Render(w, r, u.errorRenderer.RenderError(err))
		if err != nil {
			return
		}

		return
	}

	newUser, err := u.useCases.CreateAndLogin(r.Context(), payload.ToUser())
	if err != nil {
		err := render.Render(w, r, u.errorRenderer.RenderError(err))
		if err != nil {
			return
		}

		return
	}

	_ = render.Render(w, r, newLoggedUserPayload(newUser))
}

func (u *userEditPayload) Bind(_ *http.Request) error {
	v := validator.New()
	err := v.Struct(u)

	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (u *userEditPayload) ToUser() *user.User {
	return &user.User{
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Password:  u.Password,
	}
}

type userEditPayload struct {
	Name      string `json:"name" validate:"required"`
	AvatarURL string `json:"avatar_url" validate:"required"`
	Password  string `json:"password"`
}

type userCreatePayload struct {
	Name      string `json:"name"       validate:"required"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required"`
}

func (u *userCreatePayload) Bind(_ *http.Request) error {
	v := validator.New()
	err := v.Struct(u)

	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (u *userCreatePayload) ToUser() *user.User {
	return &user.User{
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Email:     u.Email,
		Password:  u.Password,
	}
}

type userPayload struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

func (u *userPayload) Bind(_ *http.Request) error {
	return nil
}

func (u *userPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newUserPayload(user *user.User) *userPayload {
	return &userPayload{
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Email:     user.Email,
	}
}

type loggedUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (l *loggedUserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newLoggedUserPayload(logged *user.Logged) *loggedUserPayload {
	return &loggedUserPayload{
		Name:  logged.Name,
		Email: logged.Email,
		Token: logged.Token,
	}
}

type userFilter struct {
	Name *string `json:"name,omitempty"`
}

func newUserList(users []*user.User) []render.Renderer {
	var list []render.Renderer

	for _, u := range users {
		list = append(list, newUserPayload(u))
	}

	return list
}
