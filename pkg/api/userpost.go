package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"tuiter.com/api/pkg/userpost"
)

func NewUserPostRouter(userPostRepository userpost.Repository) *UserPostRouter {
	return &UserPostRouter{
		userPostRepository: userPostRepository,
	}
}

type UserPostRouter struct {
	userPostRepository userpost.Repository
}

func (u UserPostRouter) Search(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(request, "page"))

	if err != nil {
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
