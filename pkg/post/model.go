package post

import (
	"gorm.io/gorm"
	"tuiter.com/api/pkg/user"
)

type Post struct {
	gorm.Model
	ParentID *int
	Message  string
	AuthorID int
	Users    []user.User `gorm:"many2many:post_likes;"`
	Likes    int
}
