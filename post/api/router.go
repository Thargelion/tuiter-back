package api

import (
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/api"
	"tuiter.com/api/kit"
	"tuiter.com/api/post/data"
)

type Router struct {
	time kit.Time
	repo data.Repository
}

func NewPostRouter(time kit.Time, repository data.Repository) *Router {
	return &Router{
		time: time,
		repo: repository,
	}
}

func (r *Router) FindAll(writer http.ResponseWriter, request *http.Request) {
	pageID := request.URL.Query().Get(string(kit.PageIDKey))
	posts, err := r.repo.FindAll(request.Context(), pageID)
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
	data := &PostPayload{}
	if err := render.Bind(request, data); err != nil {
		err := render.Render(writer, request, api.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	data.Date = r.time.Now()
	err := r.repo.Create(request.Context(), data.Post)
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
