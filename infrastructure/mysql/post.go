package mysql

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"tuiter.com/api/internal/logging"
	"tuiter.com/api/internal/post"
)

func NewPostRepository(creator *gorm.DB, logger logging.ContextualLogger) *PostRepository {
	return &PostRepository{database: creator, logger: logger}
}

type PostRepository struct {
	database *gorm.DB
	logger   logging.ContextualLogger
}

func (r *PostRepository) Create(_ context.Context, post *post.Post) error {
	res := r.database.Create(post)

	if res.Error != nil {
		return fmt.Errorf("syserror creating post %w", res.Error)
	}

	return nil
}

func (r *PostRepository) ListByPage(_ context.Context, pageID string) ([]*post.Post, error) {
	res := make([]*post.Post, 0)

	pageNumber, _ := strconv.Atoi(pageID)

	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * 100
	txResult := r.database.Limit(100).Offset(offset).Find(&res)

	if txResult.Error != nil {
		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}

func (r *PostRepository) AddLike(ctx context.Context, postID int, userID int) error {
	selectedPost, err := r.FindByID(ctx, postID)
	if err != nil {
		r.logger.Printf(ctx, "post not found when adding like %v", err)

		return fmt.Errorf("post not found when adding like %w", err)
	}
	selectedPost.Likes++

	mainTx := r.database.Begin()

	defer func() {
		if r := recover(); r != nil {
			mainTx.Rollback()
		}
	}()

	err = mainTx.Save(selectedPost).Error
	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when adding like %v", err)

		return fmt.Errorf("syserror from database when adding like %w", err)
	}

	err = mainTx.Exec("INSERT INTO post_likes (post_id, user_id) VALUES (?, ?)", postID, userID).Error
	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when registering author and tuit %v", err)

		return fmt.Errorf("syserror from database when registering author and tuit %w", err)
	}

	return mainTx.Commit().Error
}

func (r *PostRepository) RemoveLike(ctx context.Context, postID int, userID int) error {
	selectedPost, err := r.FindByID(ctx, postID)
	if err != nil {
		r.logger.Printf(ctx, "post not found when adding like %v", err)

		return fmt.Errorf("post not found when adding like %w", err)
	}
	selectedPost.Likes--

	mainTx := r.database.Begin()

	defer func() {
		if r := recover(); r != nil {
			mainTx.Rollback()
		}
	}()

	err = mainTx.Save(selectedPost).Error
	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when adding like %v", err)

		return fmt.Errorf("syserror from database when adding like %w", err)
	}

	err = mainTx.Exec(`
		DELETE FROM post_likes WHERE (post_likes.post_id = ? AND post_likes.user_id = ?)
		`, postID, userID).Error
	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when registering author and tuit %v", err)

		return fmt.Errorf("syserror from database when registering author and tuit %w", err)
	}

	return mainTx.Commit().Error
}

func (r *PostRepository) FindByID(ctx context.Context, postID int) (*post.Post, error) {
	res := &post.Post{}
	txResult := r.database.First(res, postID)

	if txResult.Error != nil {
		r.logger.Printf(ctx, "syserror from database when finding post by id %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when finding post by id %w", txResult.Error)
	}

	return res, nil
}
