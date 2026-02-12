package port

type WebhookVerifier interface {
	Verify(signature string, payload []byte) error
}
