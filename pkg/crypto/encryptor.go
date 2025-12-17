package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// KeyLength is the required length of encryption keys in bytes.
	KeyLength = 32

	// nonceLength is the length of the nonce used for encryption.
	nonceLength = 24
)

// Encryptor provides symmetric encryption for sensitive data.
// All implementations MUST be safe for concurrent use.
type Encryptor interface {
	// Encrypt encrypts plaintext and returns Base64-encoded ciphertext.
	// The output format is: Base64(nonce || encrypted_data)
	// Returns ErrEmptyPlaintext if plaintext is empty.
	Encrypt(plaintext string) (string, error)

	// Decrypt decrypts Base64-encoded ciphertext and returns plaintext.
	// Returns ErrInvalidCiphertext if input is malformed.
	// Returns ErrDecryptionFailed if authentication fails.
	Decrypt(ciphertext string) (string, error)

	// Close zeros the encryption key from memory.
	// After calling Close, the Encryptor should not be used.
	// This is recommended for long-lived server processes to minimize
	// the window during which keys are exposed in memory.
	Close() error
}

// secretboxEncryptor implements Encryptor using NaCl SecretBox.
// It is safe for concurrent use as it holds no mutable state.
type secretboxEncryptor struct {
	key [KeyLength]byte
}

// NewEncryptor creates a new Encryptor with the given 32-byte string key.
// Returns ErrInvalidKeyLength if the key is not exactly 32 bytes.
func NewEncryptor(key string) (Encryptor, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != KeyLength {
		return nil, fmt.Errorf("%w: got %d bytes", ErrInvalidKeyLength, len(keyBytes))
	}
	return newEncryptorFromBytes(keyBytes), nil
}

// NewEncryptorFromBase64 creates a new Encryptor from a Base64-encoded key.
// The decoded key must be exactly 32 bytes.
// Returns ErrInvalidKeyLength if the key is invalid or not 32 bytes.
func NewEncryptorFromBase64(encodedKey string) (Encryptor, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid base64 encoding", ErrInvalidKeyLength)
	}
	if len(keyBytes) != KeyLength {
		return nil, fmt.Errorf("%w: decoded key is %d bytes", ErrInvalidKeyLength, len(keyBytes))
	}
	return newEncryptorFromBytes(keyBytes), nil
}

// newEncryptorFromBytes creates a secretboxEncryptor from raw key bytes.
func newEncryptorFromBytes(keyBytes []byte) Encryptor {
	var key [KeyLength]byte
	copy(key[:], keyBytes)
	return &secretboxEncryptor{key: key}
}

// Encrypt encrypts plaintext using NaCl SecretBox (XSalsa20 + Poly1305).
// Output format: Base64(nonce || encrypted_data)
func (e *secretboxEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", ErrEmptyPlaintext
	}

	var nonce [nonceLength]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", fmt.Errorf("crypto: failed to generate nonce: %w", err)
	}

	encrypted := secretbox.Seal(nonce[:], []byte(plaintext), &nonce, &e.key)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts Base64-encoded ciphertext using NaCl SecretBox.
// Expected input format: Base64(nonce || encrypted_data)
func (e *secretboxEncryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", ErrInvalidCiphertext
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("%w: base64 decode failed", ErrInvalidCiphertext)
	}

	if len(data) < nonceLength+secretbox.Overhead {
		return "", fmt.Errorf("%w: ciphertext too short", ErrInvalidCiphertext)
	}

	var nonce [nonceLength]byte
	copy(nonce[:], data[:nonceLength])

	decrypted, ok := secretbox.Open(nil, data[nonceLength:], &nonce, &e.key)
	if !ok {
		return "", ErrDecryptionFailed
	}

	return string(decrypted), nil
}

// Close zeros the encryption key from memory.
// After calling Close, encrypt/decrypt operations will fail.
func (e *secretboxEncryptor) Close() error {
	for i := range e.key {
		e.key[i] = 0
	}
	return nil
}
