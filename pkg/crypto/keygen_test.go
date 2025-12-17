package crypto

import (
	"encoding/base64"
	"sync"
	"testing"
)

func Test_generateKey(t *testing.T) {
	t.Run("should generate valid base64 key", func(t *testing.T) {
		// When
		key, err := generateKey()

		// Then
		if err != nil {
			t.Fatalf("generateKey() error = %v", err)
		}
		if key == "" {
			t.Error("generateKey() returned empty key")
		}

		// Should be valid base64
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatalf("base64 decode error = %v", err)
		}

		// Should decode to 32 bytes
		if len(decoded) != KeyLength {
			t.Errorf("decoded key length = %d, want %d", len(decoded), KeyLength)
		}
	})

	t.Run("should generate unique keys", func(t *testing.T) {
		// Given
		seen := make(map[string]bool)
		iterations := 1000

		// When/Then
		for i := 0; i < iterations; i++ {
			key, err := generateKey()
			if err != nil {
				t.Fatalf("generateKey() iteration %d error = %v", i, err)
			}

			if seen[key] {
				t.Fatalf("duplicate key generated at iteration %d", i)
			}
			seen[key] = true
		}

		if len(seen) != iterations {
			t.Errorf("unique keys = %d, want %d", len(seen), iterations)
		}
	})

	t.Run("generated key should work with NewEncryptorFromBase64", func(t *testing.T) {
		// Given
		key, err := generateKey()
		if err != nil {
			t.Fatalf("generateKey() error = %v", err)
		}

		// When
		enc, err := NewEncryptorFromBase64(key)

		// Then
		if err != nil {
			t.Fatalf("NewEncryptorFromBase64() error = %v", err)
		}
		if enc == nil {
			t.Fatal("NewEncryptorFromBase64() returned nil")
		}

		// Should be usable for encryption/decryption
		plaintext := "test data"
		ciphertext, err := enc.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt() error = %v", err)
		}

		decrypted, err := enc.Decrypt(ciphertext)
		if err != nil {
			t.Fatalf("Decrypt() error = %v", err)
		}
		if decrypted != plaintext {
			t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
		}
	})

	t.Run("should be safe for concurrent use", func(t *testing.T) {
		const goroutines = 100
		var wg sync.WaitGroup
		keys := make(chan string, goroutines)
		errs := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				key, err := generateKey()
				if err != nil {
					errs <- err
					return
				}
				keys <- key
			}()
		}

		wg.Wait()
		close(keys)
		close(errs)

		// Check for errors
		for err := range errs {
			t.Errorf("generateKey() concurrent error = %v", err)
		}

		// Verify all keys are unique
		seen := make(map[string]bool)
		for key := range keys {
			if seen[key] {
				t.Error("duplicate key generated in concurrent test")
			}
			seen[key] = true
		}

		if len(seen) != goroutines {
			t.Errorf("unique keys = %d, want %d", len(seen), goroutines)
		}
	})
}

func Test_mustGenerateKey(t *testing.T) {
	t.Run("should return valid key", func(t *testing.T) {
		// When
		key := mustGenerateKey()

		// Then
		if key == "" {
			t.Error("mustGenerateKey() returned empty key")
		}

		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatalf("base64 decode error = %v", err)
		}
		if len(decoded) != KeyLength {
			t.Errorf("decoded key length = %d, want %d", len(decoded), KeyLength)
		}
	})

	t.Run("should work with NewEncryptorFromBase64", func(t *testing.T) {
		// Given
		key := mustGenerateKey()

		// When
		enc, err := NewEncryptorFromBase64(key)

		// Then
		if err != nil {
			t.Fatalf("NewEncryptorFromBase64() error = %v", err)
		}
		if enc == nil {
			t.Fatal("NewEncryptorFromBase64() returned nil")
		}
	})
}

func TestKeyLength(t *testing.T) {
	// Verify the exported constant matches expected value
	if KeyLength != 32 {
		t.Errorf("KeyLength = %d, want 32", KeyLength)
	}
}

func Benchmark_generateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = generateKey()
	}
}
