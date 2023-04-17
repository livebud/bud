package secretbox_test

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/cipher"
	"github.com/livebud/bud/package/cipher/secretbox"
)

func random(nonce []byte) error {
	copy(nonce[:], "2uVCNEZYf83uovM8aAVmNDgUdP294Ape")
	return nil
}

func generate() (secret [32]byte, err error) {
	if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
		return [32]byte{}, err
	}
	return secret, nil
}

func TestOk(t *testing.T) {
	is := is.New(t)
	secret, err := generate()
	is.NoErr(err)
	sb := secretbox.New(secret)
	sb.Random = random
	input := []byte("10")
	ciphertext, err := sb.Encrypt(input)
	is.NoErr(err)
	plaintext, err := sb.Decrypt(ciphertext)
	is.NoErr(err)
	is.Equal(string(plaintext), string(input))
}

var secretKey = [32]byte{
	0xf5, 0xaf, 0xe2, 0xcb, 0x87, 0xfb, 0x59, 0x65, 0x3d, 0xff,
	0x43, 0x56, 0x19, 0x4a, 0x22, 0x64, 0x91, 0x4a, 0x28, 0xa0,
	0x4a, 0x06, 0xb8, 0x21, 0x29, 0x42, 0xb4, 0x44, 0x55, 0xd1,
	0x13, 0x89,
}

func TestCipherText(t *testing.T) {
	is := is.New(t)
	sb := secretbox.New(secretKey)
	sb.Random = random
	input := []byte("10")
	ciphertext, err := sb.Encrypt(input)
	is.NoErr(err)
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	is.Equal("MnVWQ05FWllmODN1b3ZNOGFBVm1ORGdVB4vQmn+Qj9A8ldDGJhIIH5s3", b64)
	plaintext, err := sb.Decrypt(ciphertext)
	is.NoErr(err)
	is.Equal(string(plaintext), string(input))
}

func TestNotOk(t *testing.T) {
	is := is.New(t)
	secret, err := generate()
	is.NoErr(err)
	sb := secretbox.New(secret)
	sb.Random = random
	input := []byte("10")
	ciphertext, err := sb.Encrypt(input)
	is.NoErr(err)
	secret2, err := generate()
	is.NoErr(err)
	sb2 := secretbox.New(secret2)
	plaintext, err := sb2.Decrypt(ciphertext)
	is.True(err != nil)
	is.True(errors.Is(err, cipher.ErrDecrypting))
	is.Equal(plaintext, nil)
}
