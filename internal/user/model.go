package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}
