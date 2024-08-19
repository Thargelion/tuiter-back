package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
}

func newResponse(statusCode int, message string) *Response {
	return &Response{
		HTTPStatusCode: statusCode,
		Message:        message,
	}
}

func (response *Response) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.HTTPStatusCode)

	return nil
}
