package api

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"tuiter.com/api/post"
)

type PostPayload struct {
	*post.Post
}

func (u *PostPayload) Bind(r *http.Request) error {
	if u.Post == nil {
		return errors.New("missing required fields")
	}

	return nil
}

func (u *PostPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newPostList(posts []*post.Post) []render.Renderer {
	var list []render.Renderer
	list = []render.Renderer{}

	for _, posts := range posts {
		list = append(list, &PostPayload{posts})
	}

	return list
}
