package session

import (
	"context"
	"fmt"

	"github.com/ajg/form"
)

var ErrNotFound = fmt.Errorf("session not found")
var ErrInvalidSession = fmt.Errorf("invalid session")

func From[Session any](ctx context.Context) (session Session, err error) {
	value := ctx.Value(key)
	if value == nil {
		return session, ErrNotFound
	}
	container, ok := value.(*wrapper)
	if !ok {
		return session, ErrInvalidSession
	}
	if container.Session != nil {
		session, ok := container.Session.(Session)
		if !ok {
			return session, ErrInvalidSession
		}
		return session, nil
	}
	if err := form.DecodeString(&session, container.Raw); err != nil {
		return session, err
	}
	container.Session = session
	return session, nil
}

// - type-safe session data
// - implement swappable storage
// - allow access to the session ID
// - encrypt the cookie
// - handle parallel access (layouts, frames, etc.)
//   - user is on their own
