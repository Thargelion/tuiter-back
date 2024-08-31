package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"tuiter.com/api/internal/application/handlers"
	"tuiter.com/api/internal/application/router"
	"tuiter.com/api/internal/application/services"
	"tuiter.com/api/internal/domain/avatar"
	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/user"
	"tuiter.com/api/internal/infrastructure/mysql"
	"tuiter.com/api/pkg/logging"
)

const (
	defaultTimeoutSeconds       = 5
	defaultHeaderTimeoutSeconds = 3
)

func main() {
	// Chi
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(defaultTimeoutSeconds * time.Second))
	chiRouter.Use(handlers.RequestTagger) // Chi already has one -_-
	chiRouter.Use(middleware.Logger)
	chiRouter.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		handlers.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	addRoutes(chiRouter)

	port := os.Getenv("PORT")
	addr := ":" + port
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: defaultHeaderTimeoutSeconds * time.Second,
		Handler:           chiRouter,
	}

	printWelcomeMessage(port)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func addRoutes(chiRouter *chi.Mux) {
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DB")
	dataBase := mysql.Connect(dbUser, dbPass, dbHost, dbName)
	logger := logging.NewContextualLogger(log.Default())

	err := dataBase.AutoMigrate(&user.User{}, &tuit.Post{})
	if err != nil {
		logger.Printf(context.Background(), err.Error())
		panic("failed to migrate")
	}

	// Repositories
	userRepo := mysql.NewUserRepository(dataBase, logger)
	tuitRepo := mysql.NewTuitRepository(dataBase, logger)
	userPostRepo := mysql.NewUserPostRepository(dataBase, logger)

	// Services
	avatarUseCases := avatar.NewAvatarUseCases()
	userPostUseCases := services.NewUserPostService(tuitRepo, userPostRepo)

	// Error Handlers
	mysqlHandler := mysql.NewErrorHandler()
	errHandler := handlers.NewErrorsHandler(mysqlHandler)

	// Handlers
	userHandler := handlers.NewUserHandler(services.NewUserUseCases(userRepo, avatarUseCases), errHandler, logger)
	userPostHandler := handlers.NewUserTuitHandler(userPostUseCases, errHandler, logger)
	tuitHandler := handlers.NewTuitHandler(tuitRepo, errHandler, logger)
	likeHandler := handlers.NewLikeHandler(userPostUseCases, errHandler, logger)

	// Routers
	userRouter := router.NewUserRouter(userHandler, userPostHandler)
	tuitRouter := router.NewTuitRouter(tuitHandler)
	likeRouter := router.NewLikeRouter(likeHandler)

	chiRouter.Route("/v1", func(router chi.Router) {
		router.Route("/users", userRouter.Route)
		router.Route("/tuits", tuitRouter.Route)
		router.Route("/likes", likeRouter.Route)
	})
}

func printWelcomeMessage(port string) {
	fmt.Printf("Server running on port %s!\n", port) //nolint:forbidigo
	fmt.Print("" +                                   //nolint:forbidigo
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
