package post

import (
	"time"
	"tuiter.com/api/user/domain"
)

type Post struct {
	ID       int         `json:"id"`
	ParentId *int        `json:"parent_id"`
	Message  string      `json:"message"`
	AuthorID int         `json:"author_id"`
	Author   domain.User `json:"author"`
	Date     time.Time   `json:"date"`
}
