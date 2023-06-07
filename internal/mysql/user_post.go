package mysql

import (
	"context"
	"fmt"

	"tuiter.com/api/pkg/userpost"
)

const (
	postsPerPage = 20
)

func NewUserPostRepository(engine *GormEngine) *UserPostRepository {
	return &UserPostRepository{
		dbEngine: engine,
	}
}

type UserPostRepository struct {
	dbEngine *GormEngine
}

func (u UserPostRepository) ListByPage(_ context.Context, pageNumber int, userID int) ([]*userpost.UserPost, error) {
	var res []*userpost.UserPost

	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * postsPerPage
	txResult := u.dbEngine.Gorm.Raw(`
		SELECT p.id as id, parent_id, message, u.name as author, u.avatar_url, liked, likes, p.created_at as date 
		FROM posts as p
		    LEFT JOIN (SELECT user_id as liked, post_id FROM post_likes) pl on p.id = pl.post_id AND pl.liked = ?         
		    JOIN users u ON p.author_id = u.id
		LIMIT ?
		OFFSET ?;`,
		userID,
		postsPerPage,
		offset,
	).Scan(&res)

	if txResult.Error != nil {
		return nil, fmt.Errorf("error from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}
