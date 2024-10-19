package handlers

import (
	"net/http"
	"os"

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

func ApiValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := os.Getenv("API_KEY")
		if r.Header.Get("Api-Secret") != expected {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}
