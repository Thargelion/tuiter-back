package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type loginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (l *LoginRouter) Route(r chi.Router) {
	r.Post("/", l.handler.Login)
}

type LoginRouter struct {
	handler loginHandler
}

func NewLoginRouter(handler loginHandler) *LoginRouter {
	return &LoginRouter{handler: handler}
}
