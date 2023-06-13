package api

import (
	"net/http"

	"github.com/go-chi/render"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/pkg/post"
)

func NewPostRouter(repository post.Repository, errRenderer ErrorRenderer, logger logging.ContextualLogger) *PostRouter {
	return &PostRouter{
		repo:          repository,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type PostRouter struct {
	repo          post.Repository
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

type postPayload struct {
	commonPayload
	ParentID *int
	Message  string
	Author   *userPayload
	Likes    int
}

func newPostPayload(post *post.Post) *postPayload {
	return &postPayload{
		commonPayload: commonPayload{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
		ParentID: post.ParentID,
		Message:  post.Message,
		Author:   newUserPayload(&post.Author),
		Likes:    post.Likes,
	}
}

type createPostPayload struct {
	AuthorID int
	Message  string
}

func (c createPostPayload) Bind(_ *http.Request) error {
	if c.AuthorID == 0 || c.Message == "" {
		return errInvalidRequest
	}

	return nil
}

func (c createPostPayload) toPost() *post.Post {
	return &post.Post{
		Message:  c.Message,
		AuthorID: c.AuthorID,
	}
}

func (u *postPayload) Bind(_ *http.Request) error {
	return nil
}

func (u *postPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newPostList(posts []*post.Post) []render.Renderer {
	var list []render.Renderer

	for _, data := range posts {
		list = append(list, newPostPayload(data))
	}

	return list
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
	payload := &createPostPayload{}
	if err := render.Bind(request, payload); err != nil {
		err := render.Render(writer, request, ErrInvalidRequest(err))
		if err != nil {
			return
		}

		return
	}

	err := r.repo.Create(request.Context(), payload.toPost())

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

func (l like) Bind(_ *http.Request) error {
	if l.UserID == 0 || l.TuitID == 0 {
		return errInvalidRequest
	}

	return nil
}
