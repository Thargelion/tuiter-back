package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"tuiter.com/api/pkg/syserror"
)

var errUnknown = errors.New("unknown error")
var errTokenNotFound = errors.New("token not found")

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
	Err            error  `json:"-"`                 // low-level runtime syserror
	HTTPStatusCode int    `json:"-"`                 // handlers Response status code
	StatusText     string `json:"status"`            // handlers-level status message
	AppCode        int64  `json:"code,omitempty"`    // application-specific syserror code
	ErrorText      string `json:"message,omitempty"` // application-level syserror message, for debugging
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}

func ErrInternalServer(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error.",
		ErrorText:      err.Error(),
	}
}

func ErrNotFound(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Not Found.",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorized(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized.",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func NewErrorsHandler(handlers ...ErrorHandler) *WrapperErrorRenderer {
	return &WrapperErrorRenderer{
		errorHandler: &handlers,
	}
}

type ErrorHandler interface {
	HandleError(err error) error
}

type WrapperErrorRenderer struct {
	errorHandler *[]ErrorHandler
}

type ErrorRenderer interface {
	RenderError(err error) *ErrResponse
}

func (e *WrapperErrorRenderer) RenderError(err error) *ErrResponse {
	// Will convert specific error into a generic one
	// Will return the generic error as a Renderer
	for _, handler := range *e.errorHandler {
		err = handler.HandleError(err)
	}

	if err != nil {
		switch {
		case errors.Is(err, syserror.ErrInvalidInput):
			return ErrInvalidRequest(err)
		case errors.Is(err, syserror.ErrUnauthorized):
			return ErrUnauthorized(err)
		case errors.Is(err, syserror.ErrNotFound):
			return ErrNotFound(err)
		default:
			return ErrInternalServer(err)
		}
	}

	return ErrInternalServer(errUnknown)
}

func (e *WrapperErrorRenderer) ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
