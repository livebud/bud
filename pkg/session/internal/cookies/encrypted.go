package cookies

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/livebud/bud/internal/cipher"
)

// Encrypted cookie store
func Encrypted(c cipher.Cipher) Interface {
	return &encrypted{c}
}

// Encrypted cookie store
type encrypted struct {
	c cipher.Cipher
}

var _ Interface = (*encrypted)(nil)

func (e *encrypted) Read(r *http.Request, name string) (*http.Cookie, error) {
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
	plaintext, err := e.c.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	cookie.Value = string(plaintext)
	return cookie, err
}

func (e *encrypted) Write(w http.ResponseWriter, cookie *http.Cookie) error {
	ciphertext, err := e.c.Encrypt([]byte(cookie.Value))
	if err != nil {
		return err
	}
	cookie.Value = base64.RawURLEncoding.EncodeToString(ciphertext)
	http.SetCookie(w, cookie)
	return nil
}
