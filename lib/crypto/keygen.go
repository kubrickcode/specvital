package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// generateKey generates a new random 32-byte encryption key
// and returns it as a Base64-encoded string.
func generateKey() (string, error) {
	key := make([]byte, KeyLength)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// mustGenerateKey is like generateKey but panics if key generation fails.
// Used in tests where failure is unrecoverable.
func mustGenerateKey() string {
	key, err := generateKey()
	if err != nil {
		panic(err)
	}
	return key
}
