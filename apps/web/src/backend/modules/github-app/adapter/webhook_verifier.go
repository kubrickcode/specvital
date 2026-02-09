package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/specvital/web/src/backend/modules/github-app/domain"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

const (
	signaturePrefix        = "sha256="
	minWebhookSecretLength = 20
)

var _ port.WebhookVerifier = (*WebhookVerifier)(nil)

type WebhookVerifier struct {
	secret []byte
}

func NewWebhookVerifier(secret string) (*WebhookVerifier, error) {
	if len(secret) < minWebhookSecretLength {
		return nil, domain.ErrWeakWebhookSecret
	}
	return &WebhookVerifier{secret: []byte(secret)}, nil
}

func (v *WebhookVerifier) Verify(signature string, payload []byte) error {
	if signature == "" {
		return domain.ErrMissingSignature
	}

	if !strings.HasPrefix(signature, signaturePrefix) {
		return domain.ErrSignatureVerifyFailed
	}

	expectedMAC, err := hex.DecodeString(strings.TrimPrefix(signature, signaturePrefix))
	if err != nil {
		return domain.ErrSignatureVerifyFailed
	}

	mac := hmac.New(sha256.New, v.secret)
	mac.Write(payload)
	actualMAC := mac.Sum(nil)

	if !hmac.Equal(expectedMAC, actualMAC) {
		return domain.ErrSignatureVerifyFailed
	}

	return nil
}
