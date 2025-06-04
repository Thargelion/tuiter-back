package mysql

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"tuiter.com/api/pkg/syserror"
)

type ErrorHandler interface {
	HandleError(err error) error
}

type GormErrorHandler struct{}

func NewErrorHandler() *GormErrorHandler {
	return &GormErrorHandler{}
}

func (e GormErrorHandler) HandleError(err error) error {
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return fmt.Errorf("%w, %w", err, syserror.ErrNotFound)
		default:
			return err
		}
	}

	return nil
}
