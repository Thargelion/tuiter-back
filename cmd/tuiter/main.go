package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/internal/mysql"
	"tuiter.com/api/pkg/api"
	"tuiter.com/api/pkg/avatar"
	"tuiter.com/api/pkg/post"
	"tuiter.com/api/pkg/user"
)

func main() {
	// Chi
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(5 * time.Second))
	chiRouter.Use(api.RequestTagger) // Chi already has one -_-
	chiRouter.Use(middleware.Logger)
	chiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		api.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
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
	logger := logging.NewContextualLogger(log.Default())
	userRepo := mysql.NewUserRepository(dataBase, logger)
	avatarUseCases := avatar.NewAvatarUseCases()
	mysqlHandler := mysql.NewErrorHandler()
	postRepo := mysql.NewPostRepository(dataBase, logger)
	errHandler := api.NewErrorsHandler(mysqlHandler)
	userRouter := api.NewUserRouter(user.NewUserUseCases(userRepo, avatarUseCases), errHandler, logger)
	postRouter := api.NewPostRouter(postRepo, errHandler, logger)
	userPostRouter := api.NewUserPostRouter(mysql.NewUserPostRepository(dataBase, logger), postRepo, errHandler, logger)

	chiRouter.Route("/v1", func(router chi.Router) {
		router.Route("/users", func(r chi.Router) {
			r.Post("/", userRouter.CreateUser)
			r.Get("/{id}", userRouter.FindUserByID)
			r.Get("/", userRouter.Search)
			r.With(api.Pagination).Get("/{id}/tuits", userPostRouter.Search)
		})
		router.Route("/tuits", func(r chi.Router) {
			r.With(api.Pagination).Get("/", postRouter.Search)
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
