package post

import "net/http"

type API interface {
	Search(writer http.ResponseWriter, request *http.Request)
	CreatePost(writer http.ResponseWriter, request *http.Request)
	AddLike(writer http.ResponseWriter, request *http.Request)
}
