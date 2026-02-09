package port

type Encryptor interface {
	Decrypt(ciphertext string) (string, error)
	Encrypt(plaintext string) (string, error)
}
