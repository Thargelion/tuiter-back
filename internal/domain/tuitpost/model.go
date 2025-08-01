package tuitpost

// TuitPost represents the posts relative to a user.
// TuitPost will calculate if the user liked the tuit.
type TuitPost struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Likes     int    `json:"likes"`
	Liked     bool   `json:"liked"`
	Date      string `json:"date"`
}
