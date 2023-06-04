package mysql

import (
	"context"
	"fmt"
	"strconv"

	"tuiter.com/api/pkg/userpost"
)

func NewUserPostRepository(engine *GormEngine) *UserPostRepository {
	return &UserPostRepository{
		dbEngine: engine,
	}
}

type UserPostRepository struct {
	dbEngine *GormEngine
}

func (u UserPostRepository) ListByPage(_ context.Context, pageID string) ([]*userpost.UserPost, error) {
	var res []*userpost.UserPost

	pageNumber, _ := strconv.Atoi(pageID)

	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * 100
	txResult := u.dbEngine.Gorm.Limit(100).Offset(offset).Raw("").Find(&res)

	if txResult.Error != nil {
		return nil, fmt.Errorf("error from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}
