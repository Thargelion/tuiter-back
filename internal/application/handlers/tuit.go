package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

func NewTuitHandler(
	repository tuit.Repository,
	claimsExtractor security.TokenClaimsExtractor,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *TuitHandler {
	return &TuitHandler{
		repo:            repository,
		claimsExtractor: claimsExtractor,
		errorRenderer:   errRenderer,
		logger:          logger,
	}
}

type TuitHandler struct {
	repo            tuit.Repository
	claimsExtractor security.TokenClaimsExtractor
	errorRenderer   ErrorRenderer
	logger          logging.ContextualLogger
}

type tuitPayload struct {
	commonPayload
	ParentID *int         `json:"parent_id"`
	Message  string       `json:"message"`
	Author   *userPayload `json:"author"`
	Likes    int          `json:"likes"`
}

func newTuitPayload(post *tuit.Tuit) *tuitPayload {
	return &tuitPayload{
		commonPayload: commonPayload{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
		},
		ParentID: post.ParentID,
		Message:  post.Message,
		Author:   newUserPayload(&post.Author),
		Likes:    post.Likes,
	}
}

type createTuitPayload struct {
	AuthorID uint   `json:"author_id"`
	Message  string `json:"message"`
}

func (c createTuitPayload) Bind(_ *http.Request) error {
	if c.AuthorID == 0 || c.Message == "" {
		return errInvalidRequest
	}

	return nil
}

func (c createTuitPayload) toModel() *tuit.Tuit {
	return &tuit.Tuit{
		Message: c.Message,
		Author:  user.User{},
	}
}

func (u *tuitPayload) Bind(_ *http.Request) error {
	return nil
}

func (u *tuitPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newPostList(posts []*tuit.Tuit) []render.Renderer {
	var list []render.Renderer

	for _, data := range posts {
		list = append(list, newTuitPayload(data))
	}

	return list
}

// Search returns all the tuits from a page.
// @Summary Search tuits
// @Description Search tuits
// @Tags tuits
// @Accept json
// @Produce json
// @Param page_id query string true "Page ID"
// @Success 200 {array} tuitPayload
// @Router /tuits [get].
func (t *TuitHandler) Search(writer http.ResponseWriter, request *http.Request) {
	pageID := request.URL.Query().Get(string(PageIDKey))
	posts, err := t.repo.ListByPage(request.Context(), pageID)

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

func (t *TuitHandler) CreateTuit(w http.ResponseWriter, r *http.Request) {
	payload := &createTuitPayload{}
	if err := render.Bind(r, payload); err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			return
		}

		return
	}

	token, ok := r.Context().Value("token").(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("token not found")))
	}

	claims, err := t.claimsExtractor.ExtractClaims(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newTuit := payload.toModel()
	newTuit.Author.ID = uint(claims["sub"].(float64))

	err = t.repo.Create(r.Context(), newTuit)

	if err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	err = render.Render(w, r, newResponse(http.StatusCreated, "Tuit created"))
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
