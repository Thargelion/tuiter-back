package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"tuiter.com/api/avatar"
	"tuiter.com/api/internal/mysql"
	api2 "tuiter.com/api/pkg/api"
	"tuiter.com/api/pkg/post"
	"tuiter.com/api/pkg/user"
)

func main() {
	// Chi
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(5 * time.Second))
	chiRouter.Use(middleware.Logger)
	chiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api2.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	addRoutes(chiRouter)
	printWelcomeMessage()

	server := &http.Server{
		Addr:              ":3000",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           chiRouter,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func addRoutes(chiRouter *chi.Mux) {
	dataBase := mysql.Connect()
	err := dataBase.AutoMigrate(&user.User{}, &post.Post{})

	if err != nil {
		panic("failed to migrate")
	}

	if err != nil {
		panic("failed to load location")
	}
	// Dependencies
	userRepo := mysql.NewUserRepository(dataBase)
	avatarUseCases := avatar.NewAvatarUseCases()
	userRouter := api2.NewUserRouter(user.NewUserUseCases(userRepo, avatarUseCases))
	postRouter := api2.NewPostRouter(mysql.NewPostRepository(dataBase))
	userPostRouter := api2.NewUserPostRouter(mysql.NewUserPostRepository(dataBase))
	mockRouter := api2.NewMockRouter(dataBase)

	chiRouter.Route("/v1", func(router chi.Router) {
		router.Route("/users", func(r chi.Router) {
			r.Post("/", userRouter.CreateUser)
			r.Get("/{id}", userRouter.FindUserByID)
			r.Get("/", userRouter.Search)
			r.With(api2.Pagination).Get("/{id}/tuits", userPostRouter.Search)
		})
		router.Route("/tuits", func(r chi.Router) {
			r.With(api2.Pagination).Get("/", postRouter.Search)
			r.Post("/", postRouter.CreatePost)
		})
		router.Route("/mock", func(r chi.Router) {
			r.Post("/", mockRouter.FillMockData)
		})
	})
}

func printWelcomeMessage() {
	fmt.Print("Server running on port 3000\n") //nolint:forbidigo
	fmt.Print("" +                             //nolint:forbidigo
		"⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⣠⣤⣤⣤⣤⣤⣄⣀⡀⠄⠄⠄⠄⠄⠄⠄⠄\n" +
		"⠄⠄⠄⠄⠄⠄⠄⢀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣤⡀⠄⠄⠄⠄⠄\n" +
		"⠄⠄⠄⠄⠄⢀⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣆⠄⠄⠄⠄\n" +
		"⠄⠄⠄⠄⢠⣿⣿⣿⣿⣿⢻⣿⣿⣿⣿⣿⣿⣿⣿⣯⢻⣿⣿⣿⣿⣆⠄⠄⠄\n" +
		"⠄⠄⣼⢀⣿⣿⣿⣿⣏⡏⠄⠹⣿⣿⣿⣿⣿⣿⣿⣿⣧⢻⣿⣿⣿⣿⡆⠄⠄\n" +
		"⠄⠄⡟⣼⣿⣿⣿⣿⣿⠄⠄⠄⠈⠻⣿⣿⣿⣿⣿⣿⣿⣇⢻⣿⣿⣿⣿⠄⠄\n" +
		"⠄⢰⠃⣿⣿⠿⣿⣿⣿⠄⠄⠄⠄⠄⠄⠙⠿⣿⣿⣿⣿⣿⠄⢿⣿⣿⣿⡄⠄\n" +
		"⠄⢸⢠⣿⣿⣧⡙⣿⣿⡆⠄⠄⠄⠄⠄⠄⠄⠈⠛⢿⣿⣿⡇⠸⣿⡿⣸⡇⠄\n" +
		"⠄⠈⡆⣿⣿⣿⣿⣦⡙⠳⠄⠄⠄⠄⠄⠄⢀⣠⣤⣀⣈⠙⠃⠄⠿⢇⣿⡇⠄\n" +
		"⠄⠄⡇⢿⣿⣿⣿⣿⡇⠄⠄⠄⠄⠄⣠⣶⣿⣿⣿⣿⣿⣿⣷⣆⡀⣼⣿⡇⠄\n" +
		"⠄⠄⢹⡘⣿⣿⣿⢿⣷⡀⠄⢀⣴⣾⣟⠉⠉⠉⠉⣽⣿⣿⣿⣿⠇⢹⣿⠃⠄\n" +
		"⠄⠄⠄⢷⡘⢿⣿⣎⢻⣷⠰⣿⣿⣿⣿⣦⣀⣀⣴⣿⣿⣿⠟⢫⡾⢸⡟⠄⠄\n" +
		"⠄⠄⠄⠄⠻⣦⡙⠿⣧⠙⢷⠙⠻⠿⢿⡿⠿⠿⠛⠋⠉⠄⠂⠘⠁⠞⠄⠄⠄\n" +
		"⠄⠄⠄⠄⠄⠈⠙⠑⣠⣤⣴⡖⠄⠿⣋⣉⣉⡁⠄⢾⣦⠄⠄⠄⠄⠄⠄⠄⠄\n" +
		"⠄⠄⠄⠄⠄⠄⠄⠄⠛⠛⠋⠁⣠⣾⣿⣿⣿⣿⡆⠄⣿⠆⠄⠄⠄⠄⠄⠄⠄\n")
}
