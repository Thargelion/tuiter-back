package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"tuiter.com/api/pkg/syserror"
)

func NewBcryptErrorHandler() *BcryptErrorHandler {
	return &BcryptErrorHandler{}
}

type BcryptErrorHandler struct{}

func (bc *BcryptErrorHandler) HandleError(err error) error {
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return fmt.Errorf("%w, %s", syserror.ErrUnauthorized, "password does not match")
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return fmt.Errorf("%w, %s", syserror.ErrInvalidInput, "password is too long")
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return fmt.Errorf("%w, %s", syserror.ErrInvalidInput, "password is too short")

		default:
			return err
		}
	}

	return nil
}
