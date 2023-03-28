package cookies

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/livebud/bud/package/cipher"
)

type Store interface {
	Get(r *http.Request, name string) (*http.Cookie, error)
	Set(w http.ResponseWriter, cookie *http.Cookie) error
}

func Default() Store {
	return &defaultStore{}
}

type defaultStore struct{}

func (defaultStore) Get(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func (defaultStore) Set(w http.ResponseWriter, cookie *http.Cookie) error {
	http.SetCookie(w, cookie)
	return nil
}

// Secure cookie store
func Secure(c cipher.Cipher) Store {
	return &secure{c}
}

// secure cookie store
type secure struct {
	c cipher.Cipher
}

func (s *secure) Get(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		fmt.Println("cookie session: error getting cookie", name, err)
		return nil, err
	}
	ciphertext, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("cookie session: error decoding cookie", name, err)
		return nil, err
	}
	plaintext, err := s.c.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	cookie.Value = string(plaintext)
	return cookie, err
}

func (s *secure) Set(w http.ResponseWriter, cookie *http.Cookie) error {
	// plaintext, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	// if err != nil {
	// 	fmt.Println("cookie session: error decoding cookie", err)
	// 	return err
	// }
	ciphertext, err := s.c.Encrypt([]byte(cookie.Value))
	if err != nil {
		return err
	}
	cookie.Value = base64.RawURLEncoding.EncodeToString(ciphertext)
	http.SetCookie(w, cookie)
	return nil
}
