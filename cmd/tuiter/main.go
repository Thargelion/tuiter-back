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
	"tuiter.com/api/internal/infrastructure/mysql"
	"tuiter.com/api/pkg/instant"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/security"
)

const (
	corsMaxAge        = 300
	timeout           = 5 * time.Second
	readHeaderTimeout = 3 * time.Second
	tokenTimeoutHours = 24
	tokenTimeoutDays  = 30
)

// @title Tuiter API
// @version 1
// @description This is the API for Tuiter, a Twitter clone.
// @Schemes https
// @BasePath	/v1
// @contact.email madepietro@unlam.edu.ar.
func main() {
	// Time
	locatedTime, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	tuiterTime := instant.NewTuiterTime(locatedTime)
	// Chi
	port := os.Getenv("PORT")
	addr := ":" + port
	chiRouter := chi.NewRouter()
	// JWT
	secret := os.Getenv("JWT_SECRET")
	expiration := time.Hour * tokenTimeoutHours * tokenTimeoutDays // 30 days
	// Configure Chi
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(middleware.Timeout(timeout))
	chiRouter.Use(handlers.RequestTagger)
	chiRouter.Use(handlers.ApiValidation)
	chiRouter.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))

	// Security
	tokenValidator := security.NewJWTHandler(expiration, []byte(secret), tuiterTime)
	securityMiddleware := security.NewAuthenticatorMiddleware(tokenValidator)
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

	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DB")
	dataBase := mysql.Connect(dbUser, dbPass, dbHost, dbName)
	logger := logging.NewContextualLogger(log.Default())

	// Security

	err := dataBase.AutoMigrate(&mysql.UserEntity{}, &mysql.TuitEntity{})
	if err != nil {
		logger.Printf(context.Background(), err.Error())
		panic("failed to migrate")
	}

	// Repositories
	userRepo := mysql.NewUserRepository(dataBase, logger)
	tuitRepo := mysql.NewTuitRepository(dataBase, logger)
	userPostRepo := mysql.NewFeedRepository(dataBase, logger)

	// Services
	avatarUseCases := avatar.NewAvatarUseCases()
	userPostUseCases := services.NewUserPostService(tuitRepo, userPostRepo)
	authenticator := services.NewUserAuthenticator(userRepo, tokenValidator)
	userService := services.NewUserService(userRepo, tokenValidator, avatarUseCases)

	// Error Handlers
	mysqlHandler := mysql.NewErrorHandler()
	errHandler := handlers.NewErrorsHandler(mysqlHandler)

	// Handlers
	loginHandler := handlers.NewLogin(authenticator, errHandler)
	userHandler := handlers.NewUserHandler(userService, tokenValidator, errHandler, logger)
	userPostHandler := handlers.NewUserTuitHandler(userPostUseCases, tokenValidator, errHandler, logger)
	tuitHandler := handlers.NewTuitHandler(tuitRepo, tokenValidator, errHandler, logger)
	likeHandler := handlers.NewLikeHandler(userPostUseCases, tokenValidator, errHandler, logger)

	// Routers
	userRouter := router.NewUserRouter(userPostHandler, userHandler)
	publicUserRouter := router.NewPublicUserRouter(userHandler)
	tuitRouter := router.NewTuitRouter(tuitHandler)
	loginRouter := router.NewLoginRouter(loginHandler)
	likesRouter := router.NewLikeRouter(likeHandler)

	usersRouter := chi.NewRouter()
	usersRouter.Use(securityMiddleware.Middleware)
	usersRouter.Route("/tuits", tuitRouter.Route)
	usersRouter.Route("/me", func(router chi.Router) {
		router.Route("/", userRouter.Route)
		router.Route("/tuits/{id}/likes", likesRouter.Route)
		router.Route("/tuits", tuitRouter.Route)
	})
	chiRouter.Route("/v1", func(router chi.Router) {
		router.Route("/login", loginRouter.Route)
		router.Route("/users", publicUserRouter.Route)
		router.Mount("/", usersRouter)
	})

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           chiRouter,
	}

	printWelcomeMessage(port)

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
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
