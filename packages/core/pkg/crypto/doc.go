// Package crypto provides NaCl SecretBox encryption for secure data storage.
//
// This package implements symmetric encryption using NaCl SecretBox (XSalsa20 + Poly1305),
// which provides authenticated encryption. It is designed for encrypting sensitive data
// like OAuth tokens that need to be stored securely and decrypted later.
//
// # Security Properties
//
//   - Confidentiality: XSalsa20 stream cipher with 256-bit key
//   - Integrity: Poly1305 MAC prevents tampering
//   - Nonce: 192-bit random nonce generated per encryption (no reuse risk)
//
// # Output Format
//
// Encrypted data is returned as Base64(nonce || ciphertext), where:
//   - nonce: 24 bytes (randomly generated)
//   - ciphertext: secretbox.Seal output (plaintext + 16 bytes overhead)
//
// # Usage Example
//
//	// Create encryptor from Base64-encoded key
//	encryptor, err := crypto.NewEncryptorFromBase64(os.Getenv("ENCRYPTION_KEY"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Encrypt sensitive data
//	encrypted, err := encryptor.Encrypt(oauthToken)
//	if err != nil {
//	    return err
//	}
//
//	// Store encrypted in database...
//
//	// Later, decrypt
//	decrypted, err := encryptor.Decrypt(encrypted)
//	if err != nil {
//	    return err
//	}
//
// # Key Generation
//
// Generate a 32-byte Base64-encoded key using openssl:
//
//	openssl rand -base64 32
//
// # Thread Safety
//
// All Encryptor implementations are safe for concurrent use.
// A single Encryptor instance can be shared across goroutines.
//
// # Key Management
//
// Keys should be managed securely:
//   - Store in environment variables or secret managers (AWS Secrets Manager, HashiCorp Vault)
//   - Never commit keys to source control
//   - Use different keys for different environments
//   - Plan for key rotation (re-encrypt existing data when rotating)
//
// # Key Rotation
//
// To rotate encryption keys in production without service disruption:
//
//  1. Generate new key: openssl rand -base64 32
//
//  2. Deploy collector with BOTH keys (old for decrypt, new as fallback)
//
//  3. Deploy web with new key (encrypts with new key)
//
//  4. Run batch re-encryption of existing database tokens:
//
//     oldEnc, _ := crypto.NewEncryptorFromBase64(oldKey)
//     newEnc, _ := crypto.NewEncryptorFromBase64(newKey)
//     decrypted, _ := oldEnc.Decrypt(ciphertext)
//     newCiphertext, _ := newEnc.Encrypt(decrypted)
//
//  5. Remove old key from collector after re-encryption completes
//
// # Error Handling
//
// Wrap errors with service context for debugging distributed systems:
//
//	plaintext, err := enc.Decrypt(token)
//	if err != nil {
//	    return fmt.Errorf("decrypt token [user=%s service=%s]: %w",
//	        userID, "collector", err)
//	}
//
// # Memory Safety
//
// For long-lived server processes, call Close() when done with an Encryptor
// to zero sensitive key material from memory:
//
//	enc, _ := crypto.NewEncryptorFromBase64(key)
//	defer enc.Close()
package crypto
