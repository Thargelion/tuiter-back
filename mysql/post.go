package mysql

import (
	"context"
	"tuiter.com/api/kit"
	"tuiter.com/api/post"
)

type PostRepository struct {
	database kit.DatabaseActions
}

func (r *PostRepository) Create(ctx context.Context, post *post.Post) error {
	res := r.database.Create(post)
	return res.Error()
}

func (r *PostRepository) FindAll(ctx context.Context) ([]*post.Post, error) {
	var res []*post.Post
	txResult := r.database.Find(&res)
	return res, txResult.Error()
}

func NewPostRepository(creator kit.DatabaseActions) *PostRepository {
	return &PostRepository{database: creator}
}
