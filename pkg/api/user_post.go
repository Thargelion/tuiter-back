package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/pkg/post"
	"tuiter.com/api/pkg/userpost"
)

func NewUserPostRouter(
	userPostRepository userpost.Repository,
	postRepo post.Repository,
	errRenderer ErrorRenderer,
	logger logging.ContextualLogger,
) *UserPostRouter {
	return &UserPostRouter{
		userPostRepository: userPostRepository,
		postRepository:     postRepo,
		errorRenderer:      errRenderer,
		logger:             logger,
	}
}

type UserPostRouter struct {
	userPostRepository userpost.Repository
	postRepository     post.Repository
	errorRenderer      ErrorRenderer
	logger             logging.ContextualLogger
}

func (u *UserPostRouter) Search(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))

	if err != nil {
		u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

		page = 0
	}

	userID, err := strconv.Atoi(chi.URLParam(request, "id"))

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))

		return
	}

	userPosts, err := u.userPostRepository.ListByPage(request.Context(), page, userID)

	if err != nil {
		_ = render.Render(writer, request, ErrInvalidRequest(err))
	}

	_ = render.RenderList(writer, request, newUserPostList(userPosts))
}

func (u *UserPostRouter) AddLike(writer http.ResponseWriter, request *http.Request) {
	payload := &like{}
	if err := render.Bind(request, payload); err != nil {
		u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
		err := render.Render(writer, request, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	err := u.postRepository.AddLike(request.Context(), payload.TuitID, payload.UserID)

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	userPost, err2 := u.userPostRepository.GetByID(request.Context(), payload.UserID, payload.TuitID)

	if err2 != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	_ = render.Render(writer, request, &userPostPayload{userPost})
}

func (u *UserPostRouter) RemoveLike(writer http.ResponseWriter, request *http.Request) {
	payload := &like{}
	if err := render.Bind(request, payload); err != nil {
		u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)
		err := render.Render(writer, request, ErrInvalidRequest(err))

		if err != nil {
			return
		}

		return
	}

	err := u.postRepository.RemoveLike(request.Context(), payload.TuitID, payload.UserID)

	if err != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	userPost, err2 := u.userPostRepository.GetByID(request.Context(), payload.UserID, payload.TuitID)

	if err2 != nil {
		err := render.Render(writer, request, u.errorRenderer.RenderError(err))
		if err != nil {
			u.logger.Printf(request.Context(), "syserror rendering invalid request: %v", err)

			return
		}

		return
	}

	_ = render.Render(writer, request, &userPostPayload{userPost})
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
	return nil
}

func newUserPostList(posts []*userpost.UserPost) []render.Renderer {
	var list []render.Renderer

	for _, userPost := range posts {
		list = append(list, &userPostPayload{userPost})
	}

	return list
}
