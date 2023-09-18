package cipher

import "errors"

func Rotator(latest Interface, priors ...Interface) Interface {
	return &rotator{append([]Interface{latest}, priors...)}
}

type rotator struct {
	ciphers []Interface
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
