package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/tuitpost"
	"tuiter.com/api/pkg/query"
)

func NewUserPostService(
	tuitRepo tuit.Repository,
	userPostRepo tuitpost.Repository,
) *UserTuitService {
	return &UserTuitService{
		tuitRepository:     tuitRepo,
		userPostRepository: userPostRepo,
	}
}

type UserTuitService struct {
	tuitRepository     tuit.Repository
	userPostRepository tuitpost.Repository
}

func (u *UserTuitService) PaginateReplies(
	ctx context.Context,
	userID uint,
	tuitID uint,
	page int,
) ([]*tuitpost.TuitPost, error) {
	userTuitPage, err := u.userPostRepository.RepliesByPage(ctx, userID, tuitID, page)
	if err != nil {
		return nil, fmt.Errorf("error paginating user posts: %w", err)
	}

	return userTuitPage, nil
}

func (u *UserTuitService) Paginate(
	ctx context.Context,
	userID uint,
	page int,
	params query.Params,
) ([]*tuitpost.TuitPost, error) {
	userTuitPage, err := u.userPostRepository.SearchByPage(ctx, userID, page, params)
	if err != nil {
		return nil, fmt.Errorf("error paginating user posts: %w", err)
	}

	return userTuitPage, nil
}

func (u *UserTuitService) AddLike(
	ctx context.Context,
	userID uint,
	tuitID int,
) (*tuitpost.TuitPost, error) {
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

func (u *UserTuitService) RemoveLike(
	ctx context.Context,
	userID uint,
	tuitID int,
) (*tuitpost.TuitPost, error) {
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
