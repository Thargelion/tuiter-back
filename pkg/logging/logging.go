package logging

import (
	"context"
	"fmt"
)

const (
	uuidKey = "uuid"
)

func NewContextualLogger(logger Logger) ContextualLoggerAdapter {
	return ContextualLoggerAdapter{logger: logger}
}

type ContextualLogger interface {
	Printf(ctx context.Context, format string, v ...any)
}

type Logger interface {
	Printf(format string, v ...any)
}

type ContextualLoggerAdapter struct {
	logger Logger
}

func (c ContextualLoggerAdapter) Printf(ctx context.Context, format string, v ...any) {
	prefixedFormat := fmt.Sprint(ctx.Value(uuidKey))
	prefixedFormat += " "
	prefixedFormat += format
	c.logger.Printf(prefixedFormat, v...)
}
