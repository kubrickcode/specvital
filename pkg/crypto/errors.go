package crypto

import "errors"

// Sentinel errors for encryption operations.
var (
	// ErrInvalidKeyLength indicates the encryption key is not 32 bytes.
	ErrInvalidKeyLength = errors.New("crypto: invalid key length (expected 32 bytes)")

	// ErrDecryptionFailed indicates authentication or decryption failed.
	// This error is returned when the ciphertext cannot be authenticated,
	// typically due to wrong key or corrupted data.
	ErrDecryptionFailed = errors.New("crypto: decryption failed")

	// ErrInvalidCiphertext indicates the ciphertext is malformed.
	// This includes invalid base64 encoding or insufficient length.
	ErrInvalidCiphertext = errors.New("crypto: invalid ciphertext")

	// ErrEmptyPlaintext indicates an attempt to encrypt empty data.
	ErrEmptyPlaintext = errors.New("crypto: empty plaintext")
)
