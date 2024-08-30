package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "tuiter.com/api/cmd/tuiter/docs"
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
	corsMaxAge        = 300
	timeout           = 5 * time.Second
	readHeaderTimeout = 3 * time.Second
)

// @title Tuiter API
// @version 1
// @description This is the API for Tuiter, a Twitter clone.
// @Schemes https
// @BasePath	/v1
// @contact.email madepietro@unlam.edu.ar.
func main() {
	// Chi
	port := os.Getenv("PORT")
	addr := ":" + port
	chiRouter := chi.NewRouter()
	// Configure
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(timeout))
	chiRouter.Use(handlers.RequestTagger) // Chi already has one -_-
	chiRouter.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))
	// Add Globals
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	chiRouter.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           corsMaxAge, // Maximum value not ignored by any of major browsers
	}))
	chiRouter.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		handlers.LogWriter{ResponseWriter: w}.Write([]byte("Hello World!"))
	})
	chiRouter.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	fileServerRouter := router.NewFileServer()
	fileServerRouter.FileRoutes(chiRouter, "/files", filesDir)
	addRoutes(chiRouter)

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
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
	fmt.Printf("server running on port: %s. \n", port) //nolint:forbidigo
	fmt.Print("" +                                     //nolint:forbidigo
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
