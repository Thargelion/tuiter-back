package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"tuiter.com/api/internal/domain/userpost"
	"tuiter.com/api/pkg/logging"
)

func NewLikeHandler(liker userpost.Liker, errorRenderer ErrorRenderer, logger logging.ContextualLogger) *LikeHandler {
	return &LikeHandler{
		errorRenderer: errorRenderer,
		logger:        logger,
		liker:         liker,
	}
}

type LikeHandler struct {
	errorRenderer ErrorRenderer
	logger        logging.ContextualLogger
	liker         userpost.Liker
}

func (l *LikeHandler) AddLike(writer http.ResponseWriter, request *http.Request) {
	payload := &like{}
	if err := render.Bind(request, payload); err != nil {
		l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
		err := render.Render(writer, request, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	up, err := l.liker.AddLike(request.Context(), payload.UserID, payload.TuitID)

	if err != nil {
		_ = render.Render(writer, request, l.errorRenderer.RenderError(err))
		return
	}

	err = render.Render(writer, request, &userPostPayload{up})

	if err != nil {
		l.logger.Printf(request.Context(), "syserror rendering response: %v", err)
	}
}

func (l *LikeHandler) RemoveLike(writer http.ResponseWriter, request *http.Request) {
	payload := &like{}
	if err := render.Bind(request, payload); err != nil {
		l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
		err := render.Render(writer, request, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	up, err := l.liker.RemoveLike(request.Context(), payload.UserID, payload.TuitID)

	if err != nil {
		err := render.Render(writer, request, l.errorRenderer.RenderError(err))
		if err != nil {
			l.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	_ = render.Render(writer, request, &userPostPayload{up})
}