package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
)

type Mocker interface {
	MockData() error
}

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

type commonPayload struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type MockRouter struct {
	mocker Mocker
}

func NewMockRouter(mocker Mocker) *MockRouter {
	return &MockRouter{
		mocker: mocker,
	}
}

func (m *MockRouter) FillMockData(responseWriter http.ResponseWriter, request *http.Request) {
	err := m.mocker.MockData()
	if err != nil {
		err := render.Render(responseWriter, request, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusInternalServerError,
			StatusText:     "Internal server syserror",
			ErrorText:      err.Error(),
		})
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	err = render.Render(responseWriter, request, newResponse(http.StatusOK, "Mock data created"))

	if err != nil {
		responseWriter.WriteHeader(http.StatusOK)
	}
}

type PageKey string

const (
	// PageIDKey refers to the context key that stores the next page id.
	PageIDKey = PageKey("page")
)

// Pagination Thanks to https://github.com/jonnylangefeld/go-api/blob/v1.0.0/pkg/middelware/pagination.go !
// Pagination middleware is used to extract the next page id from the url query.
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		PageID := request.URL.Query().Get(string(PageIDKey))
		intPageID := 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.Atoi(PageID)
			if err != nil {
				_ = render.Render(responseWriter, request, ErrInvalidRequest(
					fmt.Errorf("couldn't read %s: %w esponseWriter", PageIDKey, err)),
				)

				return
			}
		}
		ctx := context.WithValue(request.Context(), PageIDKey, intPageID)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}
