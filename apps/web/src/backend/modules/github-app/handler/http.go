package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/kubrickcode/specvital/apps/web/src/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/domain/port"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/usecase"
)

const (
	headerGitHubEvent     = "X-GitHub-Event"
	headerGitHubDelivery  = "X-GitHub-Delivery"
	headerHubSignature256 = "X-Hub-Signature-256"
)

type Handler struct {
	handleWebhook *usecase.HandleWebhookUseCase
	logger        *logger.Logger
	verifier      port.WebhookVerifier
}

type HandlerConfig struct {
	HandleWebhook *usecase.HandleWebhookUseCase
	Logger        *logger.Logger
	Verifier      port.WebhookVerifier
}

var _ api.WebhookHandlers = (*Handler)(nil)

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg == nil {
		return nil, errors.New("handler config is required")
	}
	if cfg.HandleWebhook == nil {
		return nil, errors.New("HandleWebhook usecase is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Verifier == nil {
		return nil, errors.New("verifier is required")
	}

	return &Handler{
		handleWebhook: cfg.HandleWebhook,
		logger:        cfg.Logger,
		verifier:      cfg.Verifier,
	}, nil
}

func (h *Handler) HandleGitHubAppWebhookRaw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventType := r.Header.Get(headerGitHubEvent)
	deliveryID := r.Header.Get(headerGitHubDelivery)
	signature := r.Header.Get(headerHubSignature256)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(ctx, "failed to read webhook body", "error", err)
		h.respondError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	if err := h.verifier.Verify(signature, body); err != nil {
		if errors.Is(err, domain.ErrMissingSignature) || errors.Is(err, domain.ErrSignatureVerifyFailed) {
			h.logger.Warn(ctx, "webhook signature verification failed", "error", err)
			h.respondError(w, http.StatusUnauthorized, "invalid webhook signature")
			return
		}
		h.logger.Error(ctx, "webhook verification error", "error", err)
		h.respondError(w, http.StatusInternalServerError, "webhook verification failed")
		return
	}

	payload, err := parseWebhookPayload(body)
	if err != nil {
		h.logger.Error(ctx, "failed to parse webhook payload", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}

	input := usecase.HandleWebhookInput{
		Action:    payload.Action,
		EventType: eventType,
	}

	if payload.Installation != nil {
		input.InstallationID = payload.Installation.ID
		input.SuspendedAt = payload.Installation.SuspendedAt
		if payload.Installation.Account != nil {
			input.AccountID = payload.Installation.Account.ID
			input.AccountLogin = payload.Installation.Account.Login
			input.AccountType = payload.Installation.Account.Type
			input.AccountAvatarURL = payload.Installation.Account.AvatarURL
		}
	}

	output, err := h.handleWebhook.Execute(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidWebhookPayload) {
			h.logger.Warn(ctx, "invalid webhook payload", "error", err, "event", eventType)
			h.respondError(w, http.StatusBadRequest, "invalid webhook payload")
			return
		}
		h.logger.Error(ctx, "failed to handle webhook", "error", err, "event", eventType)
		h.respondError(w, http.StatusInternalServerError, "failed to process webhook")
		return
	}

	h.logger.Info(ctx, "webhook processed",
		"event", eventType,
		"action", payload.Action,
		"delivery_id", deliveryID,
		"message", output.Message,
	)

	h.respondSuccess(w, output.Message)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(api.ProblemDetail{
		Status: status,
		Title:  http.StatusText(status),
		Detail: message,
	})
}

func (h *Handler) respondSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(api.WebhookResponse{
		Success: true,
		Message: &message,
	})
}
