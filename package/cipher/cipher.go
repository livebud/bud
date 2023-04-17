package cipher

import (
	"crypto/rand"
	"errors"
	"io"
)

var ErrDecrypting = errors.New("unable to decrypt")

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

func Random(b []byte) error {
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return err
	}
	return nil
}

func Rotator(latest Cipher, priors ...Cipher) Cipher {
	return &rotator{append([]Cipher{latest}, priors...)}
}

type rotator struct {
	ciphers []Cipher
}

// Decrypt with any past or present cipher
func (r *rotator) Decrypt(ciphertext []byte) ([]byte, error) {
	for _, c := range r.ciphers {
		plaintext, err := c.Decrypt(ciphertext)
		if err != nil {
			if !errors.Is(err, ErrDecrypting) {
				return nil, err
			}
			continue
		}
		return plaintext, nil
	}
	return nil, ErrDecrypting
}

// Encrypt with the latest cipher
func (r *rotator) Encrypt(plaintext []byte) ([]byte, error) {
	return r.ciphers[0].Encrypt(plaintext)
}
