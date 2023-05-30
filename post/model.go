package post

import (
	"time"
	"tuiter.com/api/user"
)

type Post struct {
	ID       int       `json:"id"`
	ParentId *int      `json:"parent_id"`
	Message  string    `json:"message"`
	AuthorID int       `json:"author_id"`
	Author   user.User `json:"author"`
	Date     time.Time `json:"date"`
}
