package user

import "net/http"

type API interface {
	Search(writer http.ResponseWriter, request *http.Request)
	FindUserByID(writer http.ResponseWriter, request *http.Request)
	CreateUser(writer http.ResponseWriter, request *http.Request)
}
