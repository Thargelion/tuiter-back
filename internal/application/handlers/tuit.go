package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/tuitpost"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
	"tuiter.com/api/pkg/syserror"
)

func NewTuitHandler(
	repository tuit.Repository,
	tuitPostRepo tuitpost.Repository,
	userExtractor security.UserExtractor,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *TuitHandler {
	return &TuitHandler{
		repo:          repository,
		tuitPostRepo:  tuitPostRepo,
		userExtractor: userExtractor,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type TuitHandler struct {
	repo          tuit.Repository
	tuitPostRepo  tuitpost.Repository
	userExtractor security.UserExtractor
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

type tuitPayload struct {
	commonPayload
	ParentID *uint        `json:"parent_id"`
	Message  string       `json:"message"`
	Author   *userPayload `json:"author"`
	Likes    uint         `json:"likes"`
}

type tuitPostPayload struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Likes     int    `json:"likes"`
	Liked     bool   `json:"liked"`
	Date      string `json:"date"`
}

func (t *tuitPostPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newTuitPostPayload(post *tuitpost.TuitPost) *tuitPostPayload {
	return &tuitPostPayload{
		ID:        post.ID,
		Message:   post.Message,
		ParentID:  post.ParentID,
		Author:    post.Author,
		AvatarURL: post.AvatarURL,
		Likes:     post.Likes,
		Liked:     post.Liked,
		Date:      post.Date,
	}
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
	Message string `json:"message"`
}

func (c createTuitPayload) Bind(_ *http.Request) error {
	if c.Message == "" {
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
func (th *TuitHandler) Search(writer http.ResponseWriter, request *http.Request) {
	pageID := request.URL.Query().Get(string(PageIDKey))
	posts, err := th.repo.ListByPage(request.Context(), pageID)
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

func (th *TuitHandler) CreateTuit(w http.ResponseWriter, r *http.Request) {
	payload := &createTuitPayload{}
	if err := render.Bind(r, payload); err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			return
		}

		return
	}

	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errTokenNotFound))
	}

	userID, err := th.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	newTuit := payload.toModel()
	newTuit.Author.ID = userID

	err = th.repo.Create(r.Context(), newTuit)
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

func (th *TuitHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	tuitID, err := strconv.Atoi(chi.URLParam(r, "tuitID"))
	utuitID := uint(tuitID)

	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	payload := &createTuitPayload{}
	if err := render.Bind(r, payload); err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			return
		}

		return
	}

	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errTokenNotFound))
	}

	userID, err := th.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	newTuit := payload.toModel()
	newTuit.Author.ID = userID
	newTuit.ParentID = &utuitID

	err = th.repo.Create(r.Context(), newTuit)
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

// GetByID returns a single tuit by its ID.
// @Summary Get tuit by ID
// @Description Get a specific tuit by its ID
// @Tags tuits
// @Accept json
// @Produce json
// @Param tuitID path string true "Tuit ID"
// @Success 200 {object} tuitPayload
// @Router /tuits/{tuitID} [get]
func (th *TuitHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from JWT token
	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("unauthorized")))

		return
	}

	userID, err := th.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("unauthenticated")))
		return
	}
	// Get tuit ID from URL parameter
	tuitID := chi.URLParam(r, "tuitID")
	if tuitID == "" {
		_ = render.Render(
			w,
			r,
			th.errorRenderer.RenderError(
				fmt.Errorf("%w: tuit ID is required", syserror.ErrInvalidInput),
			),
		)
		return
	}

	// Convert tuitID from string to int
	id, err := strconv.Atoi(tuitID)
	if err != nil {
		_ = render.Render(
			w,
			r,
			th.errorRenderer.RenderError(
				fmt.Errorf("%w: id is not a number", syserror.ErrInvalidInput),
			),
		)
		return
	}

	// Get tuit by ID
	tuitPost, err := th.tuitPostRepo.GetByID(r.Context(), userID, id)
	if err != nil {
		th.logger.Printf(r.Context(), "error fetching tuit")
		_ = render.Render(
			w,
			r,
			th.errorRenderer.RenderError(fmt.Errorf("error retrieving tuit of id %d: %w", id, err)),
		)
		return
	}

	tuitPostPayload := newTuitPostPayload(tuitPost)

	// Render the response
	err = render.Render(w, r, tuitPostPayload)
	if err != nil {
		th.logger.Printf(r.Context(), "syserror rendering user: %v", err)

		return
	}
}
