package slots

import (
	"context"
	"fmt"
)

var ErrNotInContext = fmt.Errorf("slots: not in context")

type contextKey string

const ck contextKey = "slots"

// To returns a new context with the slots.
func (s *Slots) To(ctx context.Context) context.Context {
	return context.WithValue(ctx, ck, s)
}

// From returns the slots from the context. If the slots are not in the context,
// ErrNotInContext is returned.
func From(ctx context.Context) (*Slots, error) {
	s, ok := ctx.Value(ck).(*Slots)
	if !ok {
		return nil, ErrNotInContext
	}
	return s, nil
}
