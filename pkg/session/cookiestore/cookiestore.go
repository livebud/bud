package cookiestore

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/livebud/bud/internal/cookies"
	"github.com/livebud/bud/pkg/session"
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

func freshSession(id string) *session.State {
	return &session.State{
		ID: id,
		Payload: &session.Payload{
			Data:    map[string]any{},
			Flashes: []*session.Flash{},
		},
	}
}

func (s *Store) Load(r *http.Request, id string) (*session.State, error) {
	cookie, err := s.cs.Read(r, id)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return freshSession(id), nil
		}
		return nil, err
	}
	value, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	var payload session.Payload
	if err := s.Codec.Decode(bytes.NewBuffer(value), &payload); err != nil {
		// If there's any errors, just create a new session
		return freshSession(id), nil
	}
	if payload.Data == nil {
		payload.Data = map[string]any{}
	}
	return &session.State{
		ID:      id,
		Payload: &payload,
		Expires: cookie.Expires,
	}, nil
}

func (s *Store) Save(w http.ResponseWriter, r *http.Request, session *session.State) error {
	buffer := new(bytes.Buffer)
	if err := s.Codec.Encode(buffer, session.Payload); err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    session.ID,
		Value:   base64.RawURLEncoding.EncodeToString(buffer.Bytes()),
		Path:    "/",
		Expires: session.Expires,
	}
	return s.cs.Write(w, cookie)
}
