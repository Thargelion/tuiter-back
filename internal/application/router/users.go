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

func NewUserRouter(userPostHandler userPostHandler) *UserRouter {
	return &UserRouter{
		userPost: userPostHandler,
	}
}

func (ur *UserRouter) Route(router chi.Router) {
	router.With(handlers.Pagination).Get("/{id}/tuits", ur.userPost.Search)
}

type UserRouter struct {
	userPost userPostHandler
}

func NewPublicUserRouter(userHandler userHandler) *PublicUserRouter {
	return &PublicUserRouter{
		user: userHandler,
	}
}

func (pur *PublicUserRouter) Route(router chi.Router) {
	router.Get("/", pur.user.Search)
	router.Get("/{id}", pur.user.FindUserByID)
	router.Post("/", pur.user.CreateUser)
}

type PublicUserRouter struct {
	userPost userPostHandler
	user     userHandler
}
