package logs

import (
	"context"
	"fmt"
)

// ErrNotInContext is returned when a log is not in the context
var ErrNotInContext = fmt.Errorf("log: not in context")

type contextKey string

const logKey contextKey = "log"

// SetContext puts the log in a context and returns the new context
func ToContext(parent context.Context, log Log) context.Context {
	return context.WithValue(parent, logKey, log)
}

// From gets the log from the context. If the logger isn't in the middleware,
// we warn and discards the logs
func FromContext(ctx context.Context) (Log, error) {
	log, ok := ctx.Value(logKey).(Log)
	if !ok {
		return nil, ErrNotInContext
	}
	return log, nil
}
