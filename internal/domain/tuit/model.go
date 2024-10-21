package tuit

import (
	"time"

	"tuiter.com/api/internal/domain/user"
)

type Tuit struct {
	ID        uint
	ParentID  *uint
	Message   string
	Author    user.User
	Likes     uint
	CreatedAt time.Time
}
