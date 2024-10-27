package feed

import "context"

type Liker interface {
	AddLike
	RemoveLike
}

type AddLike interface {
	AddLike(ctx context.Context, userID uint, tuitID int) (*Feed, error)
}

type RemoveLike interface {
	RemoveLike(ctx context.Context, userID uint, tuitID int) (*Feed, error)
}

type Pager interface {
	Paginate(ctx context.Context, userID uint, page int) ([]*Feed, error)
}

type UseCases interface {
	Pager
}
