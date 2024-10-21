package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/internal/domain/feed"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

func NewLikeHandler(
	liker feed.Liker,
	userExtractor security.UserExtractor,
	errorRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *LikeHandler {
	return &LikeHandler{
		errorRenderer: errorRenderer,
		userExtractor: userExtractor,
		logger:        logger,
		liker:         liker,
	}
}

type LikeHandler struct {
	errorRenderer ErrorRenderer
	userExtractor security.UserExtractor
	logger        logging.ContextualLogger
	liker         feed.Liker
}

// AddLike godoc
// @Summary Add a like to a tuit
// @Description Add a like to a tuit
// @Tags likes
// @Accept json
// @Produce json
// @Param like body like true "Like"
// @Success 200 {object} userPostPayload
// @Router /me/tuits/{id}/likes [post].
func (l *LikeHandler) AddLike(writer http.ResponseWriter, request *http.Request) {
	tuitID, err := strconv.Atoi(chi.URLParam(request, "id"))

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(errors.New("invalid tuit id")))

		return
	}

	token, ok := request.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(writer, request, ErrInvalidRequest(errors.New("unauthorized")))

		return
	}

	userId, err := l.userExtractor.ExtractUserId(token)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userTuit, err := l.liker.AddLike(request.Context(), userId, tuitID)

	if err != nil {
		_ = render.Render(writer, request, l.errorRenderer.RenderError(err))

		return
	}

	err = render.Render(writer, request, &userPostPayload{userTuit})

	if err != nil {
		l.logger.Printf(request.Context(), "syserror rendering response: %v", err)
	}
}

// RemoveLike godoc
// @Summary Remove a like from a tuit
// @Description Remove a like from a tuit
// @Tags likes
// @Accept json
// @Produce json
// @Param like body like true "Like"
// @Success 200 {object} userPostPayload
// @Router /me/tuits/{id}/likes [delete].
func (l *LikeHandler) RemoveLike(writer http.ResponseWriter, request *http.Request) {
	tuitID, err := strconv.Atoi(chi.URLParam(request, "id"))

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(errors.New("invalid tuit id")))

		return
	}

	token, ok := request.Context().Value(security.TokenMan).(*jwt.Token)

	if !ok {
		_ = render.Render(writer, request, ErrInvalidRequest(errors.New("unauthorized")))

		return
	}

	userId, err := l.userExtractor.ExtractUserId(token)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userTuit, err := l.liker.RemoveLike(request.Context(), userId, tuitID)

	if err != nil {
		err := render.Render(writer, request, l.errorRenderer.RenderError(err))
		if err != nil {
			l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	_ = render.Render(writer, request, &userPostPayload{userTuit})
}
