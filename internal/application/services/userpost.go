package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/feed"
	"tuiter.com/api/internal/domain/tuit"
)

func NewUserPostService(tuitRepo tuit.Repository, userPostRepo feed.Repository) *UserPostService {
	return &UserPostService{
		tuitRepository:     tuitRepo,
		userPostRepository: userPostRepo,
	}
}

type UserPostService struct {
	tuitRepository     tuit.Repository
	userPostRepository feed.Repository
}

func (u *UserPostService) Paginate(ctx context.Context, userID int, page int) ([]*feed.Feed, error) {
	userTuitPage, err := u.userPostRepository.ListByPage(ctx, page, userID)

	if err != nil {
		return nil, fmt.Errorf("error paginating user posts: %w", err)
	}

	return userTuitPage, nil
}

func (u *UserPostService) AddLike(ctx context.Context, userID int, tuitID int) (*feed.Feed, error) {
	err := u.tuitRepository.AddLike(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error adding like: %w", err)
	}

	userTuit, err := u.userPostRepository.GetByID(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}

	return userTuit, nil
}

func (u *UserPostService) RemoveLike(ctx context.Context, userID int, tuitID int) (*feed.Feed, error) {
	err := u.tuitRepository.RemoveLike(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error removing like: %w", err)
	}

	userTuit, err := u.userPostRepository.GetByID(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}

	return userTuit, nil
}
