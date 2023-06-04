package post

import (
	"time"

	"gorm.io/gorm"
	"tuiter.com/api/pkg/user"
)

type Post struct {
	gorm.Model
	ParentID *int      `json:"parent_id"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date" gorm:"-"`
	AuthorID int
	Author   user.User
}
