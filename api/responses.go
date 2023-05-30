package api

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
}

func NewResponse(statusCode int, message string) render.Renderer {
	return &Response{
		HTTPStatusCode: statusCode,
		Message:        message,
	}
}

func (response *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.HTTPStatusCode)
	return nil
}
