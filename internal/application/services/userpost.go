package services

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("error adding like: %w", err)
	}

	userTuit, err := u.userPostRepository.GetByID(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}

	return userTuit, nil
}

func (u *UserPostService) RemoveLike(ctx context.Context, userID int, tuitID int) (*userpost.UserPost, error) {
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
