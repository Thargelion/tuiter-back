package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type uuidKey int

const (
	reqKey uuidKey = iota
)

func RequestTagger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier := uuid.NewString()
		ctx := context.WithValue(r.Context(), reqKey, identifier)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
