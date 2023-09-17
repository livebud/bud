package cookiestore

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/livebud/bud/pkg/session"
	"github.com/livebud/bud/pkg/session/internal/cookies"
)

// New cookie store
func New(cs cookies.Interface) *Store {
	return &Store{cs, generateRandom}
}

// Store is a cookie store
type Store struct {
	cs         cookies.Interface
	GenerateID func() (string, error)
}

var _ session.Store = (*Store)(nil)
var _ session.StoreHTTP = (*Store)(nil)

func generateRandom() (string, error) {
	// generate a random 32 byte key with valid cookie value custom alphabet
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (c *Store) loadSession(id string) (data *session.Record, err error) {
	if id == "" {
		id, err = c.GenerateID()
		if err != nil {
			return nil, err
		}
	}
	return &session.Record{
		ID:   id,
		Data: map[string]any{},
	}, nil
}

func (c *Store) LoadFrom(ctx context.Context, r *http.Request, id string) (*session.Record, error) {
	cookie, err := c.cs.Read(r, id)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return c.loadSession(id)
		}
		return nil, err
	}
	value, err := base32.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	var data session.Data
	if err := gob.NewDecoder(bytes.NewBuffer(value)).Decode(&data); err != nil {
		// If there's any errors, just create a new session
		return c.loadSession(id)
	}
	return &session.Record{
		ID:      id,
		Data:    data,
		Expires: cookie.Expires,
	}, nil
}

func (c *Store) Load(ctx context.Context, id string) (*session.Record, error) {
	return nil, errors.New("cookiestore: cannot load cookie sessions outside of the request-response lifecycle")
}

func (c *Store) SaveTo(ctx context.Context, w http.ResponseWriter, record *session.Record) error {
	value := new(bytes.Buffer)
	if err := gob.NewEncoder(value).Encode(record.Data); err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    record.ID,
		Value:   base32.StdEncoding.EncodeToString(value.Bytes()),
		Expires: record.Expires,
	}
	return c.cs.Write(w, cookie)
}

func (c *Store) Save(ctx context.Context, record *session.Record) error {
	return errors.New("cookiestore: cannot save cookie sessions outside of a request-response lifecycle")
}

func (c *Store) delete(w http.ResponseWriter, id string) error {
	cookie := &http.Cookie{
		Name:   id,
		MaxAge: -1,
	}
	return c.cs.Write(w, cookie)
}
