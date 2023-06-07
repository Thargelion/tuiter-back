package mysql

import (
	"context"
	"fmt"
	"strconv"

	"tuiter.com/api/pkg/post"
)

func NewPostRepository(creator *GormEngine) *PostRepository {
	return &PostRepository{database: creator}
}

type PostRepository struct {
	database *GormEngine
}

func (r *PostRepository) Create(_ context.Context, post *post.Post) error {
	res := r.database.Create(post)

	if res.Error() != nil {
		return fmt.Errorf("error creating post %w", res.Error())
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

	if txResult.Error() != nil {
		return nil, fmt.Errorf("error from database when listing posts by page %w", txResult.Error())
	}

	return res, nil
}
