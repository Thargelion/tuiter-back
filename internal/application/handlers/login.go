package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"tuiter.com/api/internal/domain/user"
)

func NewLogin(auth user.Authenticate, errorRenderer ErrorRenderer) *Login {
	return &Login{
		auth:          auth,
		errorRenderer: errorRenderer,
	}
}

type Login struct {
	auth          user.Authenticate
	errorRenderer ErrorRenderer
}

// Login Logs a user godoc
// @Summary Logs a user in
// @Description Logs a user in if the credentials are correct
// @Tags users
// @Accept json
// @Produce json
// @Param user body loginPayload true "User"
// @Success 200 {object} loggedUserPayload
// @Router /users [post].
func (l *Login) Login(w http.ResponseWriter, r *http.Request) {
	loginPayload := &loginPayload{}
	err := render.Bind(r, loginPayload)
	if err != nil {
		_ = render.Render(w, r, l.errorRenderer.RenderError(err))

		return
	}

	securedUser := loginPayload.toModel()

	logged, err := l.auth.Login(r.Context(), securedUser)
	if err != nil {
		renderedError := l.errorRenderer.RenderError(err)
		_ = render.Render(w, r, renderedError)

		return
	}

	_ = render.Render(w, r, newLoggedUserPayload(logged))
}

func (l *loginPayload) Bind(_ *http.Request) error {
	v := validator.New()
	err := v.Struct(l)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (l *loginPayload) toModel() *user.User {
	return &user.User{
		Email:    l.Email,
		Password: l.Password,
	}
}

type loginPayload struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
