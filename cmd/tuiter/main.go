package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
	"tuiter.com/api/kit"
	"tuiter.com/api/mysql"
	"tuiter.com/api/post"
	"tuiter.com/api/user"
)

func main() {
	db := mysql.Connect()
	err := db.AutoMigrate(&user.User{}, &post.Post{})
	if err != nil {
		panic("failed to migrate")
	}
	loc, err := time.LoadLocation("America/Buenos_Aires")
	if err != nil {
		panic("failed to load location")
	}

	// Dependencies
	tuiterTime := kit.NewTuiterTime(loc)
	userRouter := user.NewUserRouter(tuiterTime, mysql.NewUserRepository(db))
	postRouter := post.NewPostRouter(tuiterTime, mysql.NewPostRepository(db))
	mockRouter := kit.NewMockRouter(db)

	// Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		kit.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	r.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userRouter.CreateUser)
			r.Get("/{id}", userRouter.FindUserByID)
			r.Get("/", userRouter.Search)
		})
		r.Route("/posts", func(r chi.Router) {
			r.With(kit.Pagination).Get("/", postRouter.FindAll)
			r.Post("/", postRouter.CreatePost)
		})
		r.Route("/mock", func(r chi.Router) {
			r.Post("/", mockRouter.FillMockData)
		})
	})
	err = http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
