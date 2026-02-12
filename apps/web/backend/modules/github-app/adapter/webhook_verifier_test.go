package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain"
)

func TestWebhookVerifier_NewWebhookVerifier(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr error
	}{
		{
			name:    "valid secret",
			secret:  "this-is-a-valid-secret-20chars",
			wantErr: nil,
		},
		{
			name:    "short secret",
			secret:  "short",
			wantErr: domain.ErrWeakWebhookSecret,
		},
		{
			name:    "empty secret",
			secret:  "",
			wantErr: domain.ErrWeakWebhookSecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWebhookVerifier(tt.secret)
			if err != tt.wantErr {
				t.Errorf("NewWebhookVerifier() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookVerifier_Verify(t *testing.T) {
	secret := "test-webhook-secret-at-least-32-chars"
	verifier, err := NewWebhookVerifier(secret)
	if err != nil {
		t.Fatalf("NewWebhookVerifier() error = %v", err)
	}
	payload := []byte(`{"action":"opened"}`)

	validSig := generateSignature(secret, payload)

	tests := []struct {
		name      string
		payload   []byte
		signature string
		wantErr   error
	}{
		{
			name:      "valid signature",
			payload:   payload,
			signature: validSig,
			wantErr:   nil,
		},
		{
			name:      "missing signature",
			payload:   payload,
			signature: "",
			wantErr:   domain.ErrMissingSignature,
		},
		{
			name:      "wrong prefix",
			payload:   payload,
			signature: "sha1=abc",
			wantErr:   domain.ErrSignatureVerifyFailed,
		},
		{
			name:      "invalid hex",
			payload:   payload,
			signature: "sha256=notvalidhex",
			wantErr:   domain.ErrSignatureVerifyFailed,
		},
		{
			name:      "wrong signature",
			payload:   payload,
			signature: "sha256=0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:   domain.ErrSignatureVerifyFailed,
		},
		{
			name:      "modified payload",
			payload:   []byte(`{"action":"closed"}`),
			signature: validSig,
			wantErr:   domain.ErrSignatureVerifyFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifier.Verify(tt.signature, tt.payload)
			if err != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func generateSignature(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
