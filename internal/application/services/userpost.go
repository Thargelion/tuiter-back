package services

import (
	"context"

	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/userpost"
)

func NewUserPostService(tuitRepo tuit.Repository, userPostRepo userpost.Repository) *UserPostService {
	return &UserPostService{
		tuitRepository:     tuitRepo,
		userPostRepository: userPostRepo,
	}
}

type UserPostService struct {
	tuitRepository     tuit.Repository
	userPostRepository userpost.Repository
}

func (u *UserPostService) Paginate(ctx context.Context, userID int, page int) ([]*userpost.UserPost, error) {
	return u.userPostRepository.ListByPage(ctx, page, userID)
}

func (u *UserPostService) AddLike(ctx context.Context, userID int, tuitID int) (*userpost.UserPost, error) {
	err := u.tuitRepository.AddLike(ctx, userID, tuitID)

	if err != nil {
		return nil, err
	}

	return u.userPostRepository.GetByID(ctx, userID, tuitID)
}

func (u *UserPostService) RemoveLike(ctx context.Context, userID int, tuitID int) (*userpost.UserPost, error) {
	err := u.tuitRepository.RemoveLike(ctx, userID, tuitID)

	if err != nil {
		return nil, err
	}

	return u.userPostRepository.GetByID(ctx, userID, tuitID)
}
