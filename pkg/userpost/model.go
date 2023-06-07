package userpost

// UserPost represents the posts relative to a user.
// UserPost will calculate if the user liked the post.
type UserPost struct {
	ID       int
	Message  string
	ParentID int
	Author   string
	Likes    int
	Liked    bool
	Date     string
}
