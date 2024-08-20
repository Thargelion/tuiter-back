package userpost

import "context"

type Liker interface {
	AddLike
	RemoveLike
}

type AddLike interface {
	AddLike(ctx context.Context, userID int, tuitID int) (*UserPost, error)
}

type RemoveLike interface {
	RemoveLike(ctx context.Context, userID int, tuitID int) (*UserPost, error)
}

type Pager interface {
	Paginate(ctx context.Context, userID int, page int) ([]*UserPost, error)
}

type UseCases interface {
	Pager
}
