package session

import (
	"net/http"
	"time"
)

// type Storage interface {
// 	Find(id string) ([]byte, error)
// }

func New(w http.ResponseWriter, r *http.Request) *Store {
	return &Store{w, r}
}

type Store struct {
	w http.ResponseWriter
	r *http.Request
}

func (s *Store) Create(id string, payload []byte, expires time.Time) error {
	return nil
}

func (s *Store) Find(id string) ([]byte, error) {
	return nil, nil
}

func (s *Store) Delete(id string) error {
	return nil
}
