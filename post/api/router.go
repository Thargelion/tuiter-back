package api

import (
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/api"
	"tuiter.com/api/pkg"
	"tuiter.com/api/post"
)

type Router struct {
	repo post.Repository
}

func NewPostRouter(repository post.Repository) *Router {
	return &Router{
		repo: repository,
	}
}

func (r *Router) FindAll(writer http.ResponseWriter, request *http.Request) {
	pageID := request.URL.Query().Get(string(pkg.PageIDKey))
	posts, err := r.repo.ListByPage(request.Context(), pageID)
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.RenderList(writer, request, newPostList(posts))
	if err != nil {
		return
	}
}

func (r *Router) CreatePost(writer http.ResponseWriter, request *http.Request) {
	payload := &PostPayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	err := r.repo.Create(request.Context(), payload.Post)
	if err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, api.NewResponse(201, "Post created"))
	if err != nil {
		return
	}
}
