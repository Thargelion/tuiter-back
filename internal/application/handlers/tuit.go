package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/pkg/logging"
)

func NewTuitHandler(repository tuit.Repository, errRenderer ErrorRenderer, logger logging.ContextualLogger) *TuitHandler {
	return &TuitHandler{
		repo:          repository,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type TuitHandler struct {
	repo          tuit.Repository
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

type tuitPayload struct {
	commonPayload
	ParentID *int         `json:"parent_id"`
	Message  string       `json:"message"`
	Author   *userPayload `json:"author"`
	Likes    int          `json:"likes"`
}

func newTuitPayload(post *tuit.Post) *tuitPayload {
	return &tuitPayload{
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

type createTuitPayload struct {
	AuthorID int    `json:"author_id"`
	Message  string `json:"message"`
}

func (c createTuitPayload) Bind(_ *http.Request) error {
	if c.AuthorID == 0 || c.Message == "" {
		return errInvalidRequest
	}

	return nil
}

func (c createTuitPayload) toPost() *tuit.Post {
	return &tuit.Post{
		Message:  c.Message,
		AuthorID: c.AuthorID,
	}
}

func (u *tuitPayload) Bind(_ *http.Request) error {
	return nil
}

func (u *tuitPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newPostList(posts []*tuit.Post) []render.Renderer {
	var list []render.Renderer

	for _, data := range posts {
		list = append(list, newTuitPayload(data))
	}

	return list
}

func (r *TuitHandler) Search(writer http.ResponseWriter, request *http.Request) {
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

func (r *TuitHandler) CreateTuit(writer http.ResponseWriter, request *http.Request) {
	payload := &createTuitPayload{}
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
