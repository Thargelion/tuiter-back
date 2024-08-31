package tuit

import (
	"gorm.io/gorm"
	"tuiter.com/api/internal/domain/user"
)

type Post struct {
	gorm.Model
	ParentID *int
	Message  string
	AuthorID int
	Author   user.User
	Users    []user.User `gorm:"many2many:post_likes;"`
	Likes    int
}
