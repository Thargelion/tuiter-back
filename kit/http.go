package kit

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"tuiter.com/api/api"
)

type Mocker interface {
	MockData() error
}

type MockRouter struct {
	mocker Mocker
}

func NewMockRouter(mocker Mocker) *MockRouter {
	return &MockRouter{
		mocker: mocker,
	}
}

func (m *MockRouter) FillMockData(w http.ResponseWriter, r *http.Request) {
	err := m.mocker.MockData()
	if err != nil {
		err := render.Render(w, r, &api.ErrResponse{
			Err:            err,
			HTTPStatusCode: 500,
			StatusText:     "Internal server error",
			ErrorText:      err.Error(),
		})
		if err != nil {
			w.WriteHeader(500)
		}
		return
	}
	err = render.Render(w, r, api.NewResponse(200, "Mock data created"))
	if err != nil {
		w.WriteHeader(200)
	}
}

type PageKey string

const (
	// PageIDKey refers to the context key that stores the next page id
	PageIDKey = PageKey("page")
)

// Pagination Thanks to https://github.com/jonnylangefeld/go-api/blob/v1.0.0/pkg/middelware/pagination.go !
// Pagination middleware is used to extract the next page id from the url query
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PageID := r.URL.Query().Get(string(PageIDKey))
		intPageID := 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.Atoi(PageID)
			if err != nil {
				_ = render.Render(w, r, api.ErrInvalidRequest(fmt.Errorf("couldn't read %s: %w", PageIDKey, err)))
				return
			}
		}
		ctx := context.WithValue(r.Context(), PageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
