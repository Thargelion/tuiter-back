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
	"tuiter.com/api/internal/domain/tuitfeed"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

func NewTuitHandler(
	repository tuit.Repository,
	tuitFeedRepo tuitfeed.Repository,
	userExtractor security.UserExtractor,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *TuitHandler {
	return &TuitHandler{
		repo:          repository,
		tuitFeedRepo:  tuitFeedRepo,
		userExtractor: userExtractor,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type TuitHandler struct {
	repo          tuit.Repository
	tuitFeedRepo  tuitfeed.Repository
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

type tuitFeedPayload struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Likes     int    `json:"likes"`
	Liked     bool   `json:"liked"`
	Date      string `json:"date"`
}

func newTuitFeedPayload(model *tuitfeed.Model) *tuitFeedPayload {
	return &tuitFeedPayload{
		ID:        model.ID,
		Message:   model.Message,
		ParentID:  model.ParentID,
		Author:    model.Author,
		AvatarURL: model.AvatarURL,
		Likes:     model.Likes,
		Liked:     model.Liked,
		Date:      model.Date,
	}
}

func (t *tuitFeedPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
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

	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errTokenNotFound))
	}

	userID, err := t.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	newTuit := payload.toModel()
	newTuit.Author.ID = userID

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

func (t *TuitHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
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

	userID, err := t.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	newTuit := payload.toModel()
	newTuit.Author.ID = userID
	newTuit.ParentID = &utuitID

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

// GetTuitByID retrieves a specific tuit by its ID.
// @Summary Get tuit by ID
// @Description Get a tuit by its ID
// @Tags tuits
// @Accept json
// @Produce json
// @Param tuitID path int true "Tuit ID"
// @Success 200 {object} tuitFeedPayload
// @Router /tuits/{tuitID} [get]
func (t *TuitHandler) GetTuitByID(w http.ResponseWriter, r *http.Request) {
	// Extract tuitID from URL parameters
	tuitID, err := strconv.Atoi(chi.URLParam(r, "tuitID"))
	if err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			return
		}

		return
	}

	// Extract user ID from token
	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)
	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errTokenNotFound))

		return
	}

	userID, err := t.userExtractor.ExtractUserId(token)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		
		return
	}

	// Get tuit from repository
	tf, err := t.tuitFeedRepo.GetByID(r.Context(), userID, tuitID)
	if err != nil {
		t.logger.Printf(r.Context(), "error getting tuit #%d by ID: %v", tuitID, err)
		err := render.Render(w, r, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	// Check if tuit was found
	if tf == nil {
		err := render.Render(w, r, ErrNotFound(errors.New(fmt.Sprintf("tuit with id %d not found", tuitID))))
		if err != nil {
			return
		}

		return
	}

	// Return the tuit
	err = render.Render(w, r, newTuitFeedPayload(tf))
	if err != nil {
		return
	}
}
