package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"tuiter.com/api/internal/application/handlers"
)

type profileHandler interface {
	MeUser(writer http.ResponseWriter, request *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
}

type userHandler interface {
	Search(writer http.ResponseWriter, request *http.Request)
	FindUserByID(writer http.ResponseWriter, request *http.Request)
	CreateUser(writer http.ResponseWriter, request *http.Request)
}

type userPostHandler interface {
	Search(writer http.ResponseWriter, request *http.Request)
	SearchReplies(writer http.ResponseWriter, request *http.Request)
}

func NewUserRouter(userPostHandler userPostHandler, profileHandler profileHandler) *UserRouter {
	return &UserRouter{
		userPost:       userPostHandler,
		profileHandler: profileHandler,
	}
}

func (ur *UserRouter) Route(router chi.Router) {
	router.With(handlers.Pagination).Get("/feed", ur.userPost.Search)
	router.With(handlers.Pagination).Get("/feed/{tuitID}/replies", ur.userPost.SearchReplies)
	router.Get("/profile", ur.profileHandler.MeUser)
	router.Put("/profile", ur.profileHandler.UpdateProfile)
}

type UserRouter struct {
	userPost       userPostHandler
	profileHandler profileHandler
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
	user userHandler
}
