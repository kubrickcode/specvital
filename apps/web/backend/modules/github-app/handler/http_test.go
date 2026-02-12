package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubrickcode/specvital/apps/web/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/adapter"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/usecase"
)

const testWebhookSecret = "test-webhook-secret-minimum-20-chars"

func generateSignature(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestHandleGitHubAppWebhookRaw_InstallationCreated(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := map[string]interface{}{
		"action": "created",
		"installation": map[string]interface{}{
			"id": 12345,
			"account": map[string]interface{}{
				"id":         67890,
				"login":      "test-org",
				"type":       "Organization",
				"avatar_url": "https://example.com/avatar.png",
			},
		},
		"sender": map[string]interface{}{
			"id":    111,
			"login": "test-user",
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success=true")
	}
	if resp.Message != "installation created" {
		t.Errorf("expected message 'installation created', got '%s'", resp.Message)
	}
}

func TestHandleGitHubAppWebhookRaw_InstallationDeleted(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := map[string]interface{}{
		"action": "deleted",
		"installation": map[string]interface{}{
			"id": 12345,
			"account": map[string]interface{}{
				"id":    67890,
				"login": "test-org",
				"type":  "Organization",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, body))

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "installation deleted" {
		t.Errorf("expected message 'installation deleted', got '%s'", resp.Message)
	}
}

func TestHandleGitHubAppWebhookRaw_InvalidSignature(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := []byte(`{"action":"created"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalidsignature")

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandleGitHubAppWebhookRaw_MissingSignature(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := []byte(`{"action":"created"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandleGitHubAppWebhookRaw_InstallationRepositoriesAdded(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := map[string]interface{}{
		"action": "added",
		"installation": map[string]interface{}{
			"id": 12345,
			"account": map[string]interface{}{
				"id":    67890,
				"login": "test-org",
				"type":  "Organization",
			},
		},
		"repositories_added": []map[string]interface{}{
			{"id": 1, "name": "repo1", "full_name": "test-org/repo1"},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "installation_repositories")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, body))

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "repository change acknowledged" {
		t.Errorf("expected message 'repository change acknowledged', got '%s'", resp.Message)
	}
}

func TestHandleGitHubAppWebhookRaw_UnknownEvent(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := []byte(`{"action":"unknown"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "unknown_event")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, payload))

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "event type ignored" {
		t.Errorf("expected message 'event type ignored', got '%s'", resp.Message)
	}
}

func TestHandleGitHubAppWebhookRaw_InstallationSuspend(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := map[string]interface{}{
		"action": "suspend",
		"installation": map[string]interface{}{
			"id":           12345,
			"suspended_at": "2025-01-15T10:30:00Z",
			"account": map[string]interface{}{
				"id":    67890,
				"login": "test-org",
				"type":  "Organization",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, body))

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "installation suspended" {
		t.Errorf("expected message 'installation suspended', got '%s'", resp.Message)
	}
}

func TestHandleGitHubAppWebhookRaw_InstallationUnsuspend(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewHandleWebhookUseCase(repo)
	verifier, err := adapter.NewWebhookVerifier(testWebhookSecret)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	h, err := NewHandler(&HandlerConfig{
		HandleWebhook: uc,
		Logger:        logger.New(),
		Verifier:      verifier,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	payload := map[string]interface{}{
		"action": "unsuspend",
		"installation": map[string]interface{}{
			"id": 12345,
			"account": map[string]interface{}{
				"id":    67890,
				"login": "test-org",
				"type":  "Organization",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/github-app", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", generateSignature(testWebhookSecret, body))

	rr := httptest.NewRecorder()
	h.HandleGitHubAppWebhookRaw(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "installation unsuspended" {
		t.Errorf("expected message 'installation unsuspended', got '%s'", resp.Message)
	}
}
