package user

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name" gorm:"index"`
	AvatarURL string `json:"avatar_url"`
}
