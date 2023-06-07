package api

import (
	"net/http"

	"github.com/go-chi/render"
	"tuiter.com/api/pkg/post"
)

func NewPostRouter(repository post.Repository) *PostRouter {
	return &PostRouter{
		repo: repository,
	}
}

type postPayload struct {
	*post.Post
}

func (u *postPayload) Bind(_ *http.Request) error {
	if u.Post == nil {
		return errInvalidRequest
	}

	return nil
}

func (u *postPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newPostList(posts []*post.Post) []render.Renderer {
	var list []render.Renderer

	for _, data := range posts {
		list = append(list, &postPayload{data})
	}

	return list
}

type PostRouter struct {
	repo post.Repository
}

func (r *PostRouter) Search(writer http.ResponseWriter, request *http.Request) {
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

func (r *PostRouter) CreatePost(writer http.ResponseWriter, request *http.Request) {
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
