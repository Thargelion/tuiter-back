package syserror

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrInvalidInput = errors.New("invalid input")

var ErrInternal = errors.New("internal error")

var ErrUnauthorized = errors.New("unauthorized")
