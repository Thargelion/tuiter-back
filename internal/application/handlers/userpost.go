package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"tuiter.com/api/internal/domain/userpost"
	"tuiter.com/api/pkg/logging"
)

func NewUserTuitHandler(
	useCases userpost.UseCases,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *UserTuitHandler {
	return &UserTuitHandler{
		useCases:      useCases,
		errorRenderer: errRenderer,
		logger:        logger,
	}
}

type UserTuitHandler struct {
	useCases      userpost.UseCases
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
}

// Search Tuits From User godoc
// @Summary Search Users' tuits
// @Description Search Users Tuits will return a list of tuits from the user perspective. This means that the user will
// see the tuits and if they liked them or not.
// @Tags tuits
// @Param page query int false "Page"
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {array} userpost.UserPost
// @Router /users/{id}/tuits [get].
func (l *UserTuitHandler) Search(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))

	if err != nil {
		l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

		page = 0
	}

	userID, err := strconv.Atoi(chi.URLParam(request, "id"))

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userPosts, err := l.useCases.Paginate(request.Context(), page, userID)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))
	}

	_ = render.RenderList(writer, request, newUserPostList(userPosts))
}

type like struct {
	UserID int `json:"user_id"`
	TuitID int `json:"tuit_id"`
}

type userPostPayload struct {
	*userpost.UserPost
}

func (u *userPostPayload) Bind(_ *http.Request) error {
	if u.UserPost == nil {
		return errInvalidRequest
	}

	return nil
}

func (u *userPostPayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	u.Liked = u.UserPost.Liked

	return nil
}

func newUserPostList(posts []*userpost.UserPost) []render.Renderer {
	var list []render.Renderer

	for _, userPost := range posts {
		list = append(list, &userPostPayload{userPost})
	}

	return list
}
