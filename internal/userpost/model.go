package userpost

// UserPost represents the posts relative to a user.
// UserPost will calculate if the user liked the post.
type UserPost struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Likes     int    `json:"likes"`
	Liked     int    `json:"liked"`
	Date      string `json:"date"`
}