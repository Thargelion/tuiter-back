package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	api2 "tuiter.com/api/infrastructure/api"
	mysql2 "tuiter.com/api/infrastructure/mysql"
	"tuiter.com/api/internal/avatar"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/internal/post"
	user2 "tuiter.com/api/internal/user"
)

func main() {
	// Chi
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(5 * time.Second))
	chiRouter.Use(api2.RequestTagger) // Chi already has one -_-
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
	dataBase := mysql2.Connect()
	err := dataBase.AutoMigrate(&user2.User{}, &post.Post{})

	if err != nil {
		panic("failed to migrate")
	}

	if err != nil {
		panic("failed to load location")
	}
	// Dependencies
	logger := logging.NewContextualLogger(log.Default())
	userRepo := mysql2.NewUserRepository(dataBase, logger)
	avatarUseCases := avatar.NewAvatarUseCases()
	mysqlHandler := mysql2.NewErrorHandler()
	postRepo := mysql2.NewPostRepository(dataBase, logger)
	errHandler := api2.NewErrorsHandler(mysqlHandler)
	userRouter := api2.NewUserRouter(user2.NewUserUseCases(userRepo, avatarUseCases), errHandler, logger)
	postRouter := api2.NewPostRouter(postRepo, errHandler, logger)
	userPostRouter := api2.NewUserPostRouter(mysql2.NewUserPostRepository(dataBase, logger), postRepo, errHandler, logger)

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
		router.Route("/likes", func(r chi.Router) {
			r.Post("/", userPostRouter.AddLike)
		})
		router.Route("/dislikes", func(r chi.Router) {
			r.Post("/", userPostRouter.RemoveLike)
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
