package post

import (
	"gorm.io/gorm"
	"time"
	"tuiter.com/api/user"
)

type Post struct {
	gorm.Model
	ParentId *int      `json:"parent_id"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date" gorm:"-"`
	AuthorID int
	Author   user.User
}
