package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
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
	userRouter := user.NewUserRouter(mysql.NewUserRepository(db))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		kit.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	r.Post("/users", userRouter.CreateUser)
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
