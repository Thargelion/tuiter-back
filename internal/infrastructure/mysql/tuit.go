package mysql

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"tuiter.com/api/internal/domain/tuit"
	"tuiter.com/api/pkg/logging"
)

const defaultPageSize = 100

func (te *TuitEntity) TableName() string {
	return "tuits"
}

func NewTuitFromModel(t *tuit.Tuit) *TuitEntity {
	return &TuitEntity{
		Tuit:     *t,
		Message:  t.Message,
		AuthorID: t.Author.ID,
		Likes:    t.Likes,
		ParentID: t.ParentID,
	}
}

func (te *TuitEntity) ToModel() *tuit.Tuit {
	return &tuit.Tuit{
		ID:        te.Model.ID,
		ParentID:  te.ParentID,
		Message:   te.Message,
		Author:    te.Author.ToModel(),
		Likes:     te.Likes,
		CreatedAt: te.Model.CreatedAt,
	}
}

type TuitEntity struct {
	gorm.Model
	tuit.Tuit
	ParentID *uint
	Message  string
	AuthorID uint
	Author   UserEntity
	Users    []UserEntity `gorm:"many2many:tuit_likes;"`
	Likes    uint
}

func NewTuitRepository(creator *gorm.DB, logger logging.ContextualLogger) *TuitFeedRepository {
	return &TuitFeedRepository{database: creator, logger: logger}
}

type TuitFeedRepository struct {
	database *gorm.DB
	logger   logging.ContextualLogger
}

func (r *TuitFeedRepository) Create(_ context.Context, t *tuit.Tuit) error {
	res := r.database.Create(NewTuitFromModel(t))

	if res.Error != nil {
		return fmt.Errorf("syserror creating tuit %w", res.Error)
	}

	return nil
}

func (r *TuitFeedRepository) ListByPage(_ context.Context, pageID string) ([]*tuit.Tuit, error) {
	res := make([]*tuit.Tuit, 0)

	pageNumber, _ := strconv.Atoi(pageID)

	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * defaultPageSize
	txResult := r.database.Limit(defaultPageSize).Offset(offset).Find(&res)

	if txResult.Error != nil {
		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}

func (r *TuitFeedRepository) ReplyListByPage(_ context.Context, parentID uint, pageID int) ([]*tuit.Tuit, error) {
	res := make([]*tuit.Tuit, 0)

	if pageID <= 0 {
		pageID = 1
	}

	offset := (pageID - 1) * defaultPageSize
	txResult := r.database.Limit(defaultPageSize).Offset(offset).Find(&res, "parent_id = ?", parentID)

	if txResult.Error != nil {
		return nil, fmt.Errorf("syserror from database when listing posts by page %w", txResult.Error)
	}

	return res, nil
}

func (r *TuitFeedRepository) AddLike(ctx context.Context, userID uint, tuitID int) error {
	var entity TuitEntity
	err := r.database.First(&entity, tuitID).Error

	if err != nil {
		r.logger.Printf(ctx, "tuit not found when adding like %v", err)

		return fmt.Errorf("tuit not found when adding like %w", err)
	}

	entity.Likes++

	mainTx := r.database.Begin()

	defer func() {
		if r := recover(); r != nil {
			mainTx.Rollback()
		}
	}()

	err = mainTx.Save(entity).Error

	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when adding like %v", err)

		return fmt.Errorf("syserror from database when adding like %w", err)
	}

	err = mainTx.Exec("INSERT INTO tuit_likes (tuit_entity_id, user_entity_id) VALUES (?, ?)", tuitID, userID).Error
	if err != nil {
		mainTx.Rollback()
		r.logger.Printf(ctx, "syserror from database when registering author and tuit %v", err)

		return fmt.Errorf("syserror from database when registering author and tuit %w", err)
	}

	return mainTx.Commit().Error
}

func (r *TuitFeedRepository) RemoveLike(ctx context.Context, userID uint, tuitID int) error {
	var entity TuitEntity
	err := r.database.First(&entity, tuitID).Error

	if err != nil {
		r.logger.Printf(ctx, "tuit not found when adding like %v", err)

		return fmt.Errorf("tuit not found when adding like %w", err)
	}

	entity.Likes--

	txResultErr := r.database.Transaction(func(transaction *gorm.DB) error {
		err = transaction.Save(entity).Error
		if err != nil {
			r.logger.Printf(ctx, "syserror from database when adding like %v", err)

			return fmt.Errorf("syserror from database when adding like %w", err)
		}

		err = transaction.Exec(`
		DELETE FROM tuit_likes WHERE (tuit_likes.tuit_entity_id = ? AND tuit_likes.user_entity_id = ?)
		`, tuitID, userID).Error
		if err != nil {
			r.logger.Printf(ctx, "syserror from database when registering author and tuit %v", err)

			return fmt.Errorf("syserror from database when registering author and tuit %w", err)
		}

		return nil
	})

	if txResultErr != nil {
		return fmt.Errorf("syserror from database when removing like %w", txResultErr)
	}

	return nil
}

func (r *TuitFeedRepository) FindByID(ctx context.Context, tuitID int) (*tuit.Tuit, error) {
	res := &TuitEntity{}
	txResult := r.database.Preload("Author").First(res, tuitID)

	if txResult.Error != nil {
		r.logger.Printf(ctx, "syserror from database when finding tuit by id %v", txResult.Error)

		return nil, fmt.Errorf("syserror from database when finding tuit by id %w", txResult.Error)
	}

	return res.ToModel(), nil
}
