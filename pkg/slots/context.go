package slots

import (
	"context"
	"fmt"
)

var ErrNotInContext = fmt.Errorf("slots: not in context")

type contextKey string

const key contextKey = "slot"

// To returns a new context with the slots.
func ToContext(ctx context.Context, slot *Slots) context.Context {
	return context.WithValue(ctx, key, slot)
}

// From returns the slots from the context. If the slots are not in the context,
// ErrNotInContext is returned.
func FromContext(ctx context.Context) (*Slots, error) {
	s, ok := ctx.Value(key).(*Slots)
	if !ok {
		return nil, ErrNotInContext
	}
	return s, nil
}
