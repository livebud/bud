package cookiestore

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/livebud/bud/pkg/session"
	"github.com/livebud/bud/pkg/session/internal/cookies"
)

var defaultCodec = jsonCodec{}

type jsonCodec struct{}

func (jsonCodec) Encode(w io.Writer, value any) error {
	return json.NewEncoder(w).Encode(value)
}

func (jsonCodec) Decode(r io.Reader, value any) error {
	dec := json.NewDecoder(r)
	return dec.Decode(value)
}

// New cookie store
func New(cs cookies.Interface) *Store {
	return &Store{cs, defaultCodec}
}

// Store is a cookie store
type Store struct {
	cs    cookies.Interface
	Codec session.Codec
}

var _ session.Store = (*Store)(nil)

func generateRandom() (string, error) {
	// generate a random 32 byte key with valid cookie value custom alphabet
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *Store) loadSession(id string) (*session.State, error) {
	return &session.State{
		ID:   id,
		Data: map[string]any{},
	}, nil
}

func (s *Store) Load(r *http.Request, id string) (*session.State, error) {
	cookie, err := s.cs.Read(r, id)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return s.loadSession(id)
		}
		return nil, err
	}
	value, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	var data session.Data
	if err := s.Codec.Decode(bytes.NewBuffer(value), &data); err != nil {
		// If there's any errors, just create a new session
		return s.loadSession(id)
	}
	return &session.State{
		ID:      id,
		Data:    data,
		Expires: cookie.Expires,
	}, nil
}

func (s *Store) Save(w http.ResponseWriter, r *http.Request, session *session.State) error {
	buffer := new(bytes.Buffer)
	if err := s.Codec.Encode(buffer, session.Data); err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    session.ID,
		Value:   base64.RawURLEncoding.EncodeToString(buffer.Bytes()),
		Expires: session.Expires,
	}
	return s.cs.Write(w, cookie)
}
