package crypto

import (
	"encoding/base64"
	"errors"
	"strings"
	"sync"
	"testing"
)

// Test key: 32 bytes encoded in base64
const testKeyBase64 = "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY=" // "abcdefghijklmnopqrstuvwxyz123456"
const testKeyRaw = "abcdefghijklmnopqrstuvwxyz123456"

// testNewEncryptor is a test helper that creates an encryptor or fails the test.
func testNewEncryptor(t *testing.T, key string) Encryptor {
	t.Helper()
	enc, err := NewEncryptorFromBase64(key)
	if err != nil {
		t.Fatalf("NewEncryptorFromBase64() error = %v", err)
	}
	return enc
}

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr error
	}{
		{
			name:    "valid 32-byte key",
			key:     testKeyRaw,
			wantErr: nil,
		},
		{
			name:    "key too short",
			key:     "short",
			wantErr: ErrInvalidKeyLength,
		},
		{
			name:    "key too long",
			key:     "this-key-is-way-too-long-for-secretbox-encryption",
			wantErr: ErrInvalidKeyLength,
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: ErrInvalidKeyLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptor(tt.key)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("NewEncryptor() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewEncryptor() error = %v, want %v", err, tt.wantErr)
				}
				if enc != nil {
					t.Error("NewEncryptor() returned non-nil encryptor on error")
				}
			} else {
				if err != nil {
					t.Fatalf("NewEncryptor() error = %v", err)
				}
				if enc == nil {
					t.Error("NewEncryptor() returned nil")
				}
			}
		})
	}
}

func TestNewEncryptorFromBase64(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr error
	}{
		{
			name:    "valid base64 key",
			key:     testKeyBase64,
			wantErr: nil,
		},
		{
			name:    "invalid base64",
			key:     "not-valid-base64!!!",
			wantErr: ErrInvalidKeyLength,
		},
		{
			name:    "valid base64 but wrong length",
			key:     base64.StdEncoding.EncodeToString([]byte("short")),
			wantErr: ErrInvalidKeyLength,
		},
		{
			name:    "empty string",
			key:     "",
			wantErr: ErrInvalidKeyLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptorFromBase64(tt.key)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("NewEncryptorFromBase64() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewEncryptorFromBase64() error = %v, want %v", err, tt.wantErr)
				}
				if enc != nil {
					t.Error("NewEncryptorFromBase64() returned non-nil encryptor on error")
				}
			} else {
				if err != nil {
					t.Fatalf("NewEncryptorFromBase64() error = %v", err)
				}
				if enc == nil {
					t.Error("NewEncryptorFromBase64() returned nil")
				}
			}
		})
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "oauth token",
			plaintext: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
		{
			name:      "unicode text",
			plaintext: "안녕하세요 🔐 encryption test",
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("a", 10000),
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "newlines and tabs",
			plaintext: "line1\nline2\ttab",
		},
		{
			name:      "single character",
			plaintext: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			plaintext := tt.plaintext

			// When
			ciphertext, err := enc.Encrypt(plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}
			if ciphertext == "" {
				t.Error("Encrypt() returned empty ciphertext")
			}
			if ciphertext == plaintext {
				t.Error("Encrypt() returned plaintext unchanged")
			}

			// Verify base64 encoding
			if _, err := base64.StdEncoding.DecodeString(ciphertext); err != nil {
				t.Errorf("ciphertext is not valid base64: %v", err)
			}

			// Then decrypt
			decrypted, err := enc.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}
			if decrypted != plaintext {
				t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
			}
		})
	}
}

func TestEncrypt_EmptyPlaintext(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	ciphertext, err := enc.Encrypt("")
	if !errors.Is(err, ErrEmptyPlaintext) {
		t.Errorf("Encrypt() error = %v, want %v", err, ErrEmptyPlaintext)
	}
	if ciphertext != "" {
		t.Errorf("Encrypt() ciphertext = %q, want empty", ciphertext)
	}
}

func TestDecrypt_InvalidInput(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	tests := []struct {
		name       string
		ciphertext string
		wantErr    error
	}{
		{
			name:       "empty string",
			ciphertext: "",
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "invalid base64",
			ciphertext: "not-valid-base64!!!",
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "too short (less than nonce + overhead)",
			ciphertext: base64.StdEncoding.EncodeToString(make([]byte, 20)),
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "corrupted ciphertext (valid length but invalid MAC)",
			ciphertext: base64.StdEncoding.EncodeToString(make([]byte, 50)),
			wantErr:    ErrDecryptionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext, err := enc.Decrypt(tt.ciphertext)
			if err == nil {
				t.Fatalf("Decrypt() error = nil, want %v", tt.wantErr)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Decrypt() error = %v, want %v", err, tt.wantErr)
			}
			if plaintext != "" {
				t.Errorf("Decrypt() plaintext = %q, want empty", plaintext)
			}
		})
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	// Given: two different keys
	key1 := mustGenerateKey()
	key2 := mustGenerateKey()

	enc1 := testNewEncryptor(t, key1)
	enc2 := testNewEncryptor(t, key2)

	// When: encrypt with key1
	plaintext := "secret data"
	ciphertext, err := enc1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Then: decrypt with key2 should fail
	decrypted, err := enc2.Decrypt(ciphertext)
	if !errors.Is(err, ErrDecryptionFailed) {
		t.Errorf("Decrypt() with wrong key error = %v, want %v", err, ErrDecryptionFailed)
	}
	if decrypted != "" {
		t.Errorf("Decrypt() with wrong key = %q, want empty", decrypted)
	}

	// And: decrypt with key1 should succeed
	decrypted, err = enc1.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() with correct key error = %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypt_NonceRandomness(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	plaintext := "same plaintext"
	seen := make(map[string]bool)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		ciphertext, err := enc.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt() iteration %d error = %v", i, err)
		}

		// Each ciphertext should be unique (different nonce)
		if seen[ciphertext] {
			t.Fatalf("nonce collision detected at iteration %d", i)
		}
		seen[ciphertext] = true
	}

	// All ciphertexts should be unique
	if len(seen) != iterations {
		t.Errorf("unique ciphertexts = %d, want %d", len(seen), iterations)
	}
}

func TestEncryptor_Concurrent(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	const goroutines = 100
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	errChan := make(chan error, goroutines*operationsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				plaintext := "test data for concurrent access"

				// Encrypt
				ciphertext, err := enc.Encrypt(plaintext)
				if err != nil {
					errChan <- err
					return
				}

				// Decrypt
				decrypted, err := enc.Decrypt(ciphertext)
				if err != nil {
					errChan <- err
					return
				}

				if decrypted != plaintext {
					errChan <- errors.New("decrypted text does not match")
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrent operation error: %v", err)
	}
}

func TestEncryptor_CrossServiceCompatibility(t *testing.T) {
	// This test verifies that encrypted data can be decrypted by any instance
	// with the same key (simulates web -> collector scenario)

	key := mustGenerateKey()

	// Simulate web service encryptor
	webEncryptor := testNewEncryptor(t, key)

	// Simulate collector service encryptor
	collectorEncryptor := testNewEncryptor(t, key)

	// Web encrypts OAuth token
	oauthToken := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	encryptedToken, err := webEncryptor.Encrypt(oauthToken)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Collector decrypts OAuth token
	decryptedToken, err := collectorEncryptor.Decrypt(encryptedToken)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decryptedToken != oauthToken {
		t.Errorf("Decrypt() = %q, want %q", decryptedToken, oauthToken)
	}
}

func TestCiphertextFormat(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	ciphertext, err := enc.Encrypt("test")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		t.Fatalf("base64 decode error = %v", err)
	}

	// Verify minimum length: nonce (24) + plaintext (4) + overhead (16) = 44
	minLen := 24 + 4 + 16
	if len(data) < minLen {
		t.Errorf("ciphertext length = %d, want >= %d", len(data), minLen)
	}

	// First 24 bytes should be the nonce
	nonce := data[:24]
	if len(nonce) != 24 {
		t.Errorf("nonce length = %d, want 24", len(nonce))
	}

	// Rest is the encrypted data with MAC
	encryptedData := data[24:]
	minEncLen := 4 + 16 // plaintext + overhead
	if len(encryptedData) < minEncLen {
		t.Errorf("encrypted data length = %d, want >= %d", len(encryptedData), minEncLen)
	}
}

func TestEncryptor_Close(t *testing.T) {
	enc := testNewEncryptor(t, testKeyBase64)

	// Encrypt before close
	plaintext := "test data"
	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Close should zero the key
	if err := enc.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// After close, decrypt should fail (key is zeroed)
	_, err = enc.Decrypt(ciphertext)
	if !errors.Is(err, ErrDecryptionFailed) {
		t.Errorf("Decrypt() after Close() error = %v, want %v", err, ErrDecryptionFailed)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	enc, _ := NewEncryptorFromBase64(testKeyBase64)
	plaintext := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.Encrypt(plaintext)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	enc, _ := NewEncryptorFromBase64(testKeyBase64)
	ciphertext, _ := enc.Encrypt("ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.Decrypt(ciphertext)
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	enc, _ := NewEncryptorFromBase64(testKeyBase64)
	plaintext := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, _ := enc.Encrypt(plaintext)
		_, _ = enc.Decrypt(ciphertext)
	}
}
