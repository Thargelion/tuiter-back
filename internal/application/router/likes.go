package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type likeHandler interface {
	AddLike(w http.ResponseWriter, r *http.Request)
	RemoveLike(w http.ResponseWriter, r *http.Request)
}

func NewLikeRouter(likeHandler likeHandler) *LikeRouter {
	return &LikeRouter{likeHandler: likeHandler}
}

func (lr *LikeRouter) Route(r chi.Router) {
	r.Post("/", lr.likeHandler.AddLike)
	r.Delete("/", lr.likeHandler.RemoveLike)
}

type LikeRouter struct {
	likeHandler likeHandler
}
