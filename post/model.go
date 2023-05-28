package post

import "tuiter.com/api/user"

type Post struct {
	ID       int       `json:"id"`
	ParentId int       `json:"parent_id"`
	Message  string    `json:"message"`
	Author   user.User `json:"author"`
}
