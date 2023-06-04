package mysql

import (
	"context"
	"strconv"
	"tuiter.com/api/pkg"
	"tuiter.com/api/post"
)

type PostRepository struct {
	database pkg.DatabaseActions
}

func (r *PostRepository) Create(ctx context.Context, post *post.Post) error {
	res := r.database.Create(post)
	return res.Error()
}

func (r *PostRepository) ListByPage(ctx context.Context, pageId string) ([]*post.Post, error) {
	var res []*post.Post
	pageNumber, _ := strconv.Atoi(pageId)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * 100
	txResult := r.database.Limit(100).Offset(offset).Find(&res)

	return res, txResult.Error()
}

func NewPostRepository(
	creator pkg.DatabaseActions,
) *PostRepository {
	return &PostRepository{database: creator}
}
