package session

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var ErrNotFound = fmt.Errorf("session not found")
var ErrInvalidSession = fmt.Errorf("invalid session")

type Flash struct {
	Type    string
	Message string
}

type Data map[string]any

type State struct {
	ID      string    `json:"id,omitempty"`
	Data    Data      `json:"data,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

type Session struct {
	mu       sync.RWMutex
	state    *State
	modified bool
}

func (s *Session) ID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.ID
}

func (s *Session) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Data[key] = value
	s.modified = true
}

func (s *Session) SetFunc(fn func(*State)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.state)
}

func (s *Session) Get(key string) (value any) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.Data[key]
}

func (s *Session) Int(key string) (int, bool) {
	value := s.Get(key)
	switch v := value.(type) {
	case json.Number:
		i64, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return int(i64), true
	case int:
		return v, true
	case float64:
		return int(v), true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	case float32:
		return int(v), true
	default:
		return 0, false
	}
}

func (s *Session) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.state.Data[key]
	return ok
}

func (s *Session) Increment(key string) int {
	value, ok := s.Int(key)
	if !ok {
		value = 0
	}
	value++
	s.Set(key, value)
	return value
}

func (s *Session) Decrement(key string) int {
	value, ok := s.Int(key)
	if !ok {
		value = 0
	}
	value--
	s.Set(key, value)
	return value
}

func (s *Session) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Data[key]; !ok {
		return
	}
	delete(s.state.Data, key)
	s.modified = true
}

func (s *Session) Flash(kind, message string) *Flash {
	f := &Flash{
		Type:    kind,
		Message: message,
	}
	flashes, ok := s.Get("flashes").([]*Flash)
	if !ok {
		flashes = []*Flash{}
	}
	flashes = append(flashes, f)
	s.Set("flashes", flashes)
	return f
}

func (s *Session) Flashes() []*Flash {
	fmt.Println(s.state.Data)
	flashes, ok := s.Get("flashes").([]*Flash)
	if !ok {
		fmt.Println("mah...")
		return nil
	}
	s.Delete("flashes")
	return flashes
}

type Store interface {
	// Load loads the session data for the given session ID. If the session ID is
	// not found, the store should return nil data and no error.
	Load(r *http.Request, id string) (*State, error)
	Save(w http.ResponseWriter, r *http.Request, state *State) error
}

type Codec interface {
	Encode(w io.Writer, value any) error
	Decode(r io.Reader, value any) error
}

func New(store Store) *Sessions {
	return &Sessions{store: store}
}

type Sessions struct {
	store Store
}

func (s *Sessions) Load(r *http.Request, id string) (*Session, error) {
	state, err := s.store.Load(r, id)
	if err != nil {
		return nil, err
	}
	return &Session{
		state: state,
	}, nil
}

func (s *Sessions) Save(w http.ResponseWriter, r *http.Request, session *Session) error {
	if !session.modified {
		return nil
	}
	return s.store.Save(w, r, session.state)
}

const cookieKey = "sid"

type contextType string

var contextKey contextType = "session"

func (s *Sessions) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, err := s.Load(r, cookieKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(ctx, contextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
		if err := s.Save(w, r, session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func From(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(contextKey).(*Session)
	if !ok {
		return nil, ErrNotFound
	}
	return session, nil
}

// // Create a session from state
// func Create(store Store, state *State) *Session {
// 	return &Session{
// 		store: store,
// 		state: state,
// 	}
// }

// func (s *Session) Save(w http.ResponseWriter) error {
// 	return s.store.Save(w)
// }

// Fresh session
// func Fresh(id string) *Session {
// 	return &Session{
// 		id:   id,
// 		data: map[string]any{},
// 	}
// }

// func New() *Session

// // StoreHTTP is an optional interface that session stores can implement to gain
// // access to the http.Request and http.ResponseWriter.
// type StoreHTTP interface {
// 	LoadFrom(ctx context.Context, r *http.Request, id string) (*Record, error)
// 	SaveTo(ctx context.Context, w http.ResponseWriter, record *Record) error
// }

// func toStoreHTTP(store Store) StoreHTTP {
// 	if storeHTTP, ok := store.(StoreHTTP); ok {
// 		return storeHTTP
// 	}
// 	return &storeHTTP{store}
// }

// type storeHTTP struct {
// 	store Store
// }

// // LoadFrom loads from the request
// func (s *storeHTTP) LoadFrom(ctx context.Context, r *http.Request, id string) (*Record, error) {
// 	return s.store.Load(ctx, id)
// }

// // SaveTo saves to the response writer
// func (s *storeHTTP) SaveTo(ctx context.Context, w http.ResponseWriter, record *Record) error {
// 	return s.store.Save(ctx, record)
// }

// type Middleware middleware.Middleware

// func New(cookies cookies.Interface, store Store) Middleware {
// 	httpStore := toStoreHTTP(store)
// 	return middleware.Func(func(next http.Handler) http.Handler {
// 	})
// }

// func From(ctx context.Context) (session *Session, err error) {
// 	value := ctx.Value(key)
// 	if value == nil {
// 		return nil, ErrNotFound
// 	}
// 	session, ok := value.(*Session)
// 	if !ok {
// 		return nil, ErrInvalidSession
// 	}
// 	return session, nil
// }

// // - type-safe session data
// // - implement swappable storage
// // - allow access to the session ID
// // - encrypt the cookie
// // - handle parallel access (layouts, frames, etc.)
