package api

import (
	"github.com/go-chi/render"
	"net/http"
)

type response struct {
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
}

func newResponse(statusCode int, message string) render.Renderer {
	return &response{
		HTTPStatusCode: statusCode,
		Message:        message,
	}
}

func (response *response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.HTTPStatusCode)
	return nil
}
