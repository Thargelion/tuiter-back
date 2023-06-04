package post

import "net/http"

type API interface {
	FindAll(writer http.ResponseWriter, request *http.Request)
	CreatePost(writer http.ResponseWriter, request *http.Request)
}
