package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/livebud/bud/pkg/middleware"
	"github.com/livebud/bud/pkg/session/internal/cookies"
)

type Flash struct {
	Type    string
	Message string
}

type Data map[string]any

type Record struct {
	ID      string    `json:"id,omitempty"`
	Data    Data      `json:"data,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

type Session struct {
	mu       sync.RWMutex
	record   *Record
	modified bool
}

func (s *Session) ID() string {
	return s.record.ID
}

func (s *Session) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.record.Data[key] = value
	s.modified = true
}

func (s *Session) Get(key string) (value any) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.record.Data[key]
}

func (s *Session) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.record.Data[key]; !ok {
		return
	}
	delete(s.record.Data, key)
	s.modified = true
}

func (s *Session) Flash(kind, message string) *Flash {
	f := &Flash{
		Type:    kind,
		Message: message,
	}
	// s.data.Flashes = append(s.data.Flashes, f)
	s.modified = true
	return f
}

type Store interface {
	// Load loads the session data for the given session ID. If the session ID is
	// not found, the store should return nil data and no error.
	Load(ctx context.Context, id string) (*Record, error)
	Save(ctx context.Context, record *Record) error
}

// StoreHTTP is an optional interface that session stores can implement to gain
// access to the http.Request and http.ResponseWriter.
type StoreHTTP interface {
	LoadFrom(ctx context.Context, r *http.Request, id string) (*Record, error)
	SaveTo(ctx context.Context, w http.ResponseWriter, record *Record) error
}

func toStoreHTTP(store Store) StoreHTTP {
	if storeHTTP, ok := store.(StoreHTTP); ok {
		return storeHTTP
	}
	return &storeHTTP{store}
}

type storeHTTP struct {
	store Store
}

// LoadFrom loads from the request
func (s *storeHTTP) LoadFrom(ctx context.Context, r *http.Request, id string) (*Record, error) {
	return s.store.Load(ctx, id)
}

// SaveTo saves to the response writer
func (s *storeHTTP) SaveTo(ctx context.Context, w http.ResponseWriter, record *Record) error {
	return s.store.Save(ctx, record)
}

type Middleware middleware.Middleware

func New(cookies cookies.Interface, store Store) Middleware {
	httpStore := toStoreHTTP(store)
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sessionID := ""
			cookie, err := cookies.Read(r, cookieName)
			if nil == err {
				sessionID = cookie.Value
			} else if !errors.Is(err, http.ErrNoCookie) {
				// TODO: figure out what to do
				panic(err)
			}
			sessionData, err := httpStore.LoadFrom(ctx, r, sessionID)
			if err != nil {
				panic(err)
			}
			session := &Session{
				record: sessionData,
			}
			ctx = context.WithValue(ctx, key, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			if !session.modified {
				return
			}
			if err := httpStore.SaveTo(ctx, w, session.record); err != nil {
				panic(err)
			}
			if err := cookies.Write(w, &http.Cookie{
				Name:     cookieName,
				Value:    session.record.ID,
				Expires:  session.record.Expires,
				HttpOnly: true,
			}); err != nil {
				panic(err)
			}
		})
	})
}

const cookieName = "sid"

type contextKey string

var key contextKey = "session"

var ErrNotFound = fmt.Errorf("session not found")
var ErrInvalidSession = fmt.Errorf("invalid session")

func From(ctx context.Context) (session *Session, err error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, ErrNotFound
	}
	session, ok := value.(*Session)
	if !ok {
		return nil, ErrInvalidSession
	}
	return session, nil
}

// // - type-safe session data
// // - implement swappable storage
// // - allow access to the session ID
// // - encrypt the cookie
// // - handle parallel access (layouts, frames, etc.)
