package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"tuiter.com/api/internal/application/handlers"
)

type tuitHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	CreateTuit(w http.ResponseWriter, r *http.Request)
	CreateReply(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}

func NewTuitRouter(tuitHandler tuitHandler) *TuitRouter {
	return &TuitRouter{tuitHandler: tuitHandler}
}

type TuitRouter struct {
	tuitHandler tuitHandler
}

func (tr *TuitRouter) Route(r chi.Router) {
	r.With(handlers.Pagination).Get("/", tr.tuitHandler.Search)
	r.Post("/", tr.tuitHandler.CreateTuit)
	r.Post("/{tuitID}/replies", tr.tuitHandler.CreateReply)
	r.Get("/{tuitID}", tr.tuitHandler.GetByID)
}
