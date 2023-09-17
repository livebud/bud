package cipher

import (
	"errors"
)

var ErrDecrypting = errors.New("unable to decrypt")

type Interface interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
