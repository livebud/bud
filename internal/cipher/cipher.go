package cipher

import "errors"

var ErrDecrypting = errors.New("unable to decrypt")

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

func Plain() Cipher {
	return plain{}
}

type plain struct{}

func (plain) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (plain) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}

func Encrypted(key [32]byte) Cipher {
	return encrypted{key: key}
}

type encrypted struct {
	key [32]byte
}

func (e encrypted) Encrypt(plaintext []byte) ([]byte, error) {
	// return encrypt(e.key, plaintext)
	return nil, ErrDecrypting
}

func (e encrypted) Decrypt(ciphertext []byte) ([]byte, error) {
	// return decrypt(e.key, ciphertext)
	return nil, ErrDecrypting
}
