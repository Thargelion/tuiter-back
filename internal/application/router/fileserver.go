package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func NewFileServer() *FileServer {
	return &FileServer{}
}

type FileServer struct{}

// FileRoutes conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (f *FileServer) FileRoutes(router chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileRoutes does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	path += "*"

	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
