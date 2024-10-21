package feed

// Feed represents the posts relative to a user.
// Feed will calculate if the user liked the tuit.
type Feed struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Likes     int    `json:"likes"`
	Liked     bool   `json:"liked"`
	Date      string `json:"date"`
}
