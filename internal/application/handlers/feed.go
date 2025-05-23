package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/internal/domain/tuitfeed"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/query"
	"tuiter.com/api/pkg/security"
)

func NewUserTuitHandler(
	useCases tuitfeed.UseCases,
	claimsExtractor security.TokenClaimsExtractor,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *UserTuitHandler {
	return &UserTuitHandler{
		useCases:        useCases,
		claimsExtractor: claimsExtractor,
		errorRenderer:   errRenderer,
		logger:          logger,
	}
}

type UserTuitHandler struct {
	useCases        tuitfeed.UseCases
	claimsExtractor security.TokenClaimsExtractor
	errorRenderer   ErrorRenderer
	logger          logging.ContextualLogger
}

// Search Tuits From User godoc
// @Summary Search Users' tuits
// @Description Search Users Tuits will return a list of tuits from the user perspective. This means that the user will
// see the tuits and if they liked them or not.
// @Tags tuits
// @Param page query int false "Page"
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {array} tuitfeed.tuitfeed
// @Router /me/feed [get].
func (l *UserTuitHandler) Search(w http.ResponseWriter, r *http.Request) {
	test := r.URL.Query().Get("page")

	l.logger.Printf(r.Context(), "page: %s", test)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		l.logger.Printf(r.Context(), "syserror rendering invalid r: %v", err)

		page = 0
	}

	token, ok := r.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	claims, err := l.claimsExtractor.ExtractClaims(token)

	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))

		return
	}

	userID := uint(claims["sub"].(float64))

	userTuits, err := l.useCases.Paginate(r.Context(), userID, page, query.FromURLQuery(r.URL.Query()))

	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
	}

	_ = render.RenderList(w, r, newuserTuitList(userTuits))
}

// SearchReplies Tuits From User godoc
// @Summary Search Users' tuits
// @Description Search Users Tuits will return a list of tuits from the user perspective. This means that the user will
// see the tuits and if they liked them or not.
// @Tags tuits
// @Param page query int false "Page"
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {array} tuitfeed.tuitfeed
// @Router /me/tuits/{tuitID}/replies [get].
func (l *UserTuitHandler) SearchReplies(writer http.ResponseWriter, request *http.Request) {
	tuitID, err := strconv.Atoi(chi.URLParam(request, "tuitID"))
	uTuitID := uint(tuitID)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	page, err := strconv.Atoi(request.URL.Query().Get("page"))

	if err != nil {
		l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

		page = 0
	}

	token, ok := request.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	claims, err := l.claimsExtractor.ExtractClaims(token)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userID := uint(claims["sub"].(float64))

	userTuits, err := l.useCases.PaginateReplies(request.Context(), userID, uTuitID, page)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))
	}

	_ = render.RenderList(writer, request, newuserTuitList(userTuits))
}

type userTuitPayload struct {
	*tuitfeed.Model
}

func (u *userTuitPayload) Bind(_ *http.Request) error {
	if u.Model == nil {
		return errInvalidRequest
	}

	return nil
}

func (u *userTuitPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	u.Liked = u.Model.Liked

	return nil
}

func newuserTuitList(posts []*tuitfeed.Model) []render.Renderer {
	list := []render.Renderer{}

	for _, userTuit := range posts {
		list = append(list, &userTuitPayload{userTuit})
	}

	return list
}
