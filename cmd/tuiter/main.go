package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"tuiter.com/api/kit"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		kit.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
