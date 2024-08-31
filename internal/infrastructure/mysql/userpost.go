package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"tuiter.com/api/internal/domain/userpost"
	"tuiter.com/api/pkg/logging"
)

const (
	postsPerPage              = 20
	projectedPostPartialQuery = `
	SELECT p.id as id, parent_id, message, u.name as author, u.avatar_url, pl.user_id IS NOT NULL as liked , 
	       likes, p.created_at as date
		FROM posts as p
		    LEFT JOIN (SELECT user_id, post_id FROM post_likes) pl on p.id = pl.post_id AND pl.user_id = ?         
		    JOIN users u ON p.author_id = u.id
	`
	projectedPostByIDQuery = projectedPostPartialQuery + "WHERE p.id = ?"
)

func NewUserPostRepository(engine *gorm.DB, logger logging.ContextualLogger) *UserPostRepository {
	return &UserPostRepository{
		dbEngine: engine,
		logger:   logger,
	}
}

type UserPostRepository struct {
	dbEngine *gorm.DB
	logger   logging.ContextualLogger
}

func (u UserPostRepository) GetByID(ctx context.Context, userID int, postID int) (*userpost.UserPost, error) {
	var res *userpost.UserPost
	txResult := u.dbEngine.Raw(
		projectedPostByIDQuery+";",
		userID,
		postID).Scan(&res)

	if txResult.Error != nil {
		u.logger.Printf(ctx, "syserror from database when listing posts by page %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, txResult.Error
}

func (u UserPostRepository) ListByPage(ctx context.Context, userID int, page int) ([]*userpost.UserPost, error) {
	var res []*userpost.UserPost

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * postsPerPage
	txResult := u.selectUserPosts(userID, offset).Scan(&res)

	if txResult.Error != nil {
		u.logger.Printf(ctx, "syserror from database when listing posts by page %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}

func (u UserPostRepository) selectUserPosts(userID int, offset int) *gorm.DB {
	return u.dbEngine.Raw(
		paginatedQuery(projectedPostPartialQuery),
		userID,
		postsPerPage,
		offset,
	)
}
