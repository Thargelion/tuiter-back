package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"tuiter.com/api/internal/application/handlers"
)

type userHandler interface {
	Search(writer http.ResponseWriter, request *http.Request)
	FindUserByID(writer http.ResponseWriter, request *http.Request)
	CreateUser(writer http.ResponseWriter, request *http.Request)
}

type userPostHandler interface {
	Search(writer http.ResponseWriter, request *http.Request)
}

func NewUserRouter(userHandler userHandler, userPostHandler userPostHandler) *UserRouter {
	return &UserRouter{
		user:     userHandler,
		userPost: userPostHandler,
	}
}

func (ur *UserRouter) Route(router chi.Router) {
	router.Get("/", ur.user.Search)
	router.Get("/{id}", ur.user.FindUserByID)
	router.Post("/", ur.user.CreateUser)
	router.With(handlers.Pagination).Get("/{id}/tuits", ur.userPost.Search)
}

type UserRouter struct {
	user     userHandler
	userPost userPostHandler
}
