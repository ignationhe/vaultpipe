package envfile

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

// DefaultEncryptOptions returns sensible defaults for encryption.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		Placeholder: "[encrypted]",
	}
}

// EncryptOptions controls encryption and decryption behaviour.
type EncryptOptions struct {
	// Placeholder is used when displaying encrypted values in dry-run output.
	Placeholder string
}

// deriveKey produces a 32-byte AES-256 key from a passphrase via SHA-256.
func deriveKey(passphrase string) []byte {
	sum := sha256.Sum256([]byte(passphrase))
	return sum[:]
}

// Encrypt encrypts every value in secrets using AES-256-GCM and returns a new
// map whose values are base64-encoded ciphertext strings.
func Encrypt(secrets map[string]string, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("encrypt: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encrypt: create GCM: %w", err)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, fmt.Errorf("encrypt: generate nonce for %q: %w", k, err)
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		out[k] = base64.StdEncoding.EncodeToString(ciphertext)
	}
	return out, nil
}

// Decrypt decrypts every value in secrets that was previously produced by
// Encrypt and returns a map of plaintext values.
func Decrypt(secrets map[string]string, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("decrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt: create GCM: %w", err)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, fmt.Errorf("decrypt: base64 decode %q: %w", k, err)
		}
		if len(data) < gcm.NonceSize() {
			return nil, fmt.Errorf("decrypt: ciphertext too short for key %q", k)
		}
		nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
		plain, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt: open %q: %w", k, err)
		}
		out[k] = string(plain)
	}
	return out, nil
}
