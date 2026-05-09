// Package encrypt provides symmetric encryption and decryption of .env
// values using AES-GCM so that sensitive entries can be stored safely at rest.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidCiphertext is returned when a ciphertext cannot be decoded or
// authenticated.
var ErrInvalidCiphertext = errors.New("encrypt: invalid or tampered ciphertext")

// deriveKey produces a 32-byte AES-256 key from an arbitrary passphrase via
// SHA-256.  For production use, prefer a proper KDF such as Argon2id.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// EncryptValue encrypts plaintext with AES-256-GCM using a key derived from
// passphrase and returns a base64url-encoded ciphertext (nonce || ciphertext).
func EncryptValue(passphrase, plaintext string) (string, error) {
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", fmt.Errorf("encrypt: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: new gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: generate nonce: %w", err)
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(sealed), nil
}

// DecryptValue reverses EncryptValue.  It returns ErrInvalidCiphertext when
// the data cannot be decoded or the authentication tag does not match.
func DecryptValue(passphrase, encoded string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", fmt.Errorf("encrypt: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: new gcm: %w", err)
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", ErrInvalidCiphertext
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plain), nil
}

// EncryptMap returns a new map where every value whose key is listed in keys
// has been replaced with its encrypted form.  Keys not listed are copied as-is.
func EncryptMap(passphrase string, src map[string]string, keys []string) (map[string]string, error) {
	sensitive := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		sensitive[k] = struct{}{}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if _, ok := sensitive[k]; ok {
			enc, err := EncryptValue(passphrase, v)
			if err != nil {
				return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
			}
			out[k] = enc
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// DecryptMap is the inverse of EncryptMap.
func DecryptMap(passphrase string, src map[string]string, keys []string) (map[string]string, error) {
	sensitive := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		sensitive[k] = struct{}{}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if _, ok := sensitive[k]; ok {
			dec, err := DecryptValue(passphrase, v)
			if err != nil {
				return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
			}
			out[k] = dec
		} else {
			out[k] = v
		}
	}
	return out, nil
}
