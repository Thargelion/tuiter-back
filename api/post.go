package api

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/post"
)

func NewPostRouter(repository post.Repository) post.Api {
	return &router{
		repo: repository,
	}
}

type postPayload struct {
	*post.Post
}

func (u *postPayload) Bind(r *http.Request) error {
	if u.Post == nil {
		return errors.New("missing required fields")
	}

	return nil
}

func (u *postPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newPostList(posts []*post.Post) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, posts := range posts {
		list = append(list, &postPayload{posts})
	}

	return list
}

type router struct {
	repo post.Repository
}

func (r *router) FindAll(writer http.ResponseWriter, request *http.Request) {
	pageID := request.URL.Query().Get(string(PageIDKey))
	posts, err := r.repo.ListByPage(request.Context(), pageID)
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
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

func (r *router) CreatePost(writer http.ResponseWriter, request *http.Request) {
	payload := &postPayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}
	err := r.repo.Create(request.Context(), payload.Post)
	if err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, newResponse(201, "Post created"))
	if err != nil {
		return
	}
}
