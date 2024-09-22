package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/internal/domain/userpost"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

func NewUserTuitHandler(
	useCases userpost.UseCases,
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
	useCases        userpost.UseCases
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
// @Success 200 {array} userpost.UserPost
// @Router /me/feed [get].
func (l *UserTuitHandler) Search(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))

	if err != nil {
		l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

		page = 0
	}

	token, ok := request.Context().Value("token").(*jwt.Token)

	if !ok {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	claims, err := l.claimsExtractor.ExtractClaims(token)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userID := int(claims["sub"].(float64))

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
	list := []render.Renderer{}

	for _, userPost := range posts {
		list = append(list, &userPostPayload{userPost})
	}

	return list
}
