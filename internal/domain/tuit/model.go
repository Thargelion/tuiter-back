package tuit

import (
	"tuiter.com/api/internal/domain/user"
)

type Post struct {
	ParentID *int
	Message  string
	AuthorID int
	Author   user.User
	Users    []user.User
	Likes    int
}
