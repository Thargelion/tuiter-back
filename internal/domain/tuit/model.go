package tuit

import (
	"time"

	"tuiter.com/api/internal/domain/user"
)

type Tuit struct {
	ID        uint
	ParentID  *int
	Message   string
	Author    user.User
	Likes     int
	CreatedAt time.Time
}
