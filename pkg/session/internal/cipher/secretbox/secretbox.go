package secretbox

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/livebud/bud/pkg/session/internal/cipher"
	"golang.org/x/crypto/nacl/secretbox"
)

func defaultRandom(b []byte) error {
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return err
	}
	return nil
}

func New(secret [32]byte) *Cipher {
	return &Cipher{
		secret: &secret,
		Random: defaultRandom,
	}
}

type Cipher struct {
	secret *[32]byte
	Random func([]byte) error
}

var _ cipher.Interface = (*Cipher)(nil)

// Encrypt the plaintext using the secret key.
func (c *Cipher) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if err := c.Random(nonce[:]); err != nil {
		return nil, err
	}
	// This encrypts "hello world" and appends the result to the nonce.
	ciphertext = secretbox.Seal(nonce[:], plaintext, &nonce, c.secret)
	return ciphertext, nil
}

// Decrypt the ciphertext using the secret key.
func (c *Cipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	plaintext, ok := secretbox.Open(nil, ciphertext[24:], &nonce, c.secret)
	if !ok {
		return nil, fmt.Errorf("secretbox: %w", cipher.ErrDecrypting)
	}
	return plaintext, nil
}
