package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"tuiter.com/api/internal/domain/tuitpost"
	"tuiter.com/api/pkg/logging"
	"tuiter.com/api/pkg/query"
)

const (
	postsPerPage              = 20
	projectedPostPartialQuery = `
	SELECT t.id as id, parent_id, message, u.name as author, u.avatar_url, pl.user_entity_id IS NOT NULL as liked , 
	       likes, t.created_at as date 
		FROM tuits as t
		    LEFT JOIN (SELECT user_entity_id, tuit_entity_id FROM tuit_likes) pl 
		on t.id = pl.tuit_entity_id AND pl.user_entity_id = ?         
		    JOIN users u ON t.author_id = u.id
	`
	projectedPostByIDQuery = projectedPostPartialQuery + "WHERE t.id = ?"
)

func NewFeedRepository(engine *gorm.DB, logger logging.ContextualLogger) *FeedRepository {
	return &FeedRepository{
		dbEngine: engine,
		logger:   logger,
	}
}

type FeedRepository struct {
	dbEngine *gorm.DB
	logger   logging.ContextualLogger
}

func (u FeedRepository) GetByID(ctx context.Context, userID uint, postID int) (*tuitpost.TuitPost, error) {
	var res *tuitpost.TuitPost
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

func (u FeedRepository) RepliesByPage(ctx context.Context, userID uint, parentID uint, page int) ([]*tuitpost.TuitPost, error) {
	var res []*tuitpost.TuitPost

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * postsPerPage
	q := query.Builder(
		projectedPostPartialQuery,
	).Where(
		"t.parent_id = ?",
	).OrderBy(
		"t.created_at desc",
	).Paginated(
		postsPerPage, offset,
	)

	txResult := u.dbEngine.Raw(q.String(), userID, parentID).Scan(&res)

	if txResult.Error != nil {
		u.logger.Printf(ctx, "syserror from database when listing posts by page %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}

func (u FeedRepository) SearchByPage(ctx context.Context, userID uint, page int, params query.Params) ([]*tuitpost.TuitPost, error) {
	var res []*tuitpost.TuitPost

	if page <= 0 {
		page = 1
	}

	onlyParents := params.Contains("only_parents")

	offset := (page - 1) * postsPerPage
	q := query.Builder(
		projectedPostPartialQuery,
	)
	if onlyParents {
		q = q.Where(
			"t.parent_id IS NULL",
		)
	}
	q = q.OrderBy(
		"t.created_at desc",
	).Paginated(
		postsPerPage,
		offset,
	)
	txResult := u.dbEngine.Raw(
		string(q),
		userID,
	).Scan(&res)

	if txResult.Error != nil {
		u.logger.Printf(ctx, "syserror from database when listing posts by page %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}
