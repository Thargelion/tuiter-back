package services

import (
	"context"
	"fmt"

	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/internal/domain/tuitfeed"
	"tuiter.com/api/pkg/query"
)

func NewuserTuitService(tuitRepo tuit.Repository, userTuitRepo tuitfeed.Repository) *UserTuitService {
	return &UserTuitService{
		tuitRepository:     tuitRepo,
		tuitFeedRepository: userTuitRepo,
	}
}

type UserTuitService struct {
	tuitRepository     tuit.Repository
	tuitFeedRepository tuitfeed.Repository
}

func (u *UserTuitService) PaginateReplies(
	ctx context.Context,
	userID uint,
	tuitID uint,
	page int,
) ([]*tuitfeed.Model, error) {
	userTuitPage, err := u.tuitFeedRepository.RepliesByPage(ctx, userID, tuitID, page)

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
) ([]*tuitfeed.Model, error) {
	userTuitPage, err := u.tuitFeedRepository.SearchByPage(ctx, userID, page, params)

	if err != nil {
		return nil, fmt.Errorf("error paginating user posts: %w", err)
	}

	return userTuitPage, nil
}

func (u *UserTuitService) AddLike(ctx context.Context, userID uint, tuitID int) (*tuitfeed.Model, error) {
	err := u.tuitRepository.AddLike(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error adding like: %w", err)
	}

	userTuit, err := u.tuitFeedRepository.GetByID(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}

	return userTuit, nil
}

func (u *UserTuitService) RemoveLike(ctx context.Context, userID uint, tuitID int) (*tuitfeed.Model, error) {
	err := u.tuitRepository.RemoveLike(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error removing like: %w", err)
	}

	userTuit, err := u.tuitFeedRepository.GetByID(ctx, userID, tuitID)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}

	return userTuit, nil
}

func (u *UserTuitService) GetByID(ctx context.Context, userID uint, tuitID int) (*tuitfeed.Model, error) {
	post, err := u.tuitFeedRepository.GetByID(ctx, userID, tuitID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user post: %w", err)
	}
	return post, nil
}
