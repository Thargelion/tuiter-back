package post

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/rest"
)

type Router struct {
	repo Repository
}

func NewPostRouter(repository Repository) *Router {
	return &Router{repository}
}

func (r *Router) FindAll(writer http.ResponseWriter, request *http.Request) {
	posts, err := r.repo.FindAll(request.Context())
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
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
	data := &Payload{}
	if err := render.Bind(request, data); err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err := r.repo.Create(request.Context(), data.Post)
	if err != nil {
		err := render.Render(writer, request, rest.ErrInvalidRequest(err))
		if err != nil {
			return
		}
		return
	}

	err = render.Render(writer, request, rest.NewResponse(201, "Post created"))
	if err != nil {
		return
	}
}

type Payload struct {
	*Post
}

func (u *Payload) Bind(r *http.Request) error {
	if u.Post == nil {
		return errors.New("missing required User fields")
	}

	return nil
}

func (u *Payload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newPostList(posts []*Post) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, posts := range posts {
		list = append(list, &Payload{posts})
	}

	return list
}
