package userpost

import "net/http"

type API interface {
	Search(writer http.ResponseWriter, request *http.Request)
	Create(writer http.ResponseWriter, request *http.Request)
}
