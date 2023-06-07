package api

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type LogWriter struct {
	http.ResponseWriter
}

func (w LogWriter) Write(p []byte) {
	_, err := w.ResponseWriter.Write(p)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // api Response status code
	StatusText     string `json:"status"`          // api-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func ErrInvalidRequest(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}
