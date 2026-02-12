package usecase

import (
	"context"
	"time"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain/port"
)

type HandleWebhookInput struct {
	Action           string
	EventType        string
	InstallationID   int64
	AccountID        int64
	AccountLogin     string
	AccountType      string
	AccountAvatarURL *string
	SuspendedAt      *string
}

type HandleWebhookOutput struct {
	Message string
}

type HandleWebhookUseCase struct {
	repo port.InstallationRepository
}

func NewHandleWebhookUseCase(repo port.InstallationRepository) *HandleWebhookUseCase {
	return &HandleWebhookUseCase{repo: repo}
}

func (uc *HandleWebhookUseCase) Execute(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	switch input.EventType {
	case "installation":
		return uc.handleInstallation(ctx, input)
	case "installation_repositories":
		return uc.handleInstallationRepositories(ctx, input)
	default:
		return &HandleWebhookOutput{Message: "event type ignored"}, nil
	}
}

func (uc *HandleWebhookUseCase) handleInstallation(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	switch input.Action {
	case "created":
		return uc.handleInstallationCreated(ctx, input)
	case "deleted":
		return uc.handleInstallationDeleted(ctx, input)
	case "suspend":
		return uc.handleInstallationSuspend(ctx, input)
	case "unsuspend":
		return uc.handleInstallationUnsuspend(ctx, input)
	default:
		return &HandleWebhookOutput{Message: "installation action ignored"}, nil
	}
}

func (uc *HandleWebhookUseCase) handleInstallationCreated(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	if input.InstallationID <= 0 || input.AccountID <= 0 {
		return nil, domain.ErrInvalidWebhookPayload
	}

	accountType := entity.AccountTypeUser
	if input.AccountType == "Organization" {
		accountType = entity.AccountTypeOrganization
	}

	installation := &entity.Installation{
		InstallationID:   input.InstallationID,
		AccountID:        input.AccountID,
		AccountLogin:     input.AccountLogin,
		AccountType:      accountType,
		AccountAvatarURL: input.AccountAvatarURL,
	}

	if err := uc.repo.Upsert(ctx, installation); err != nil {
		return nil, err
	}

	return &HandleWebhookOutput{Message: "installation created"}, nil
}

func (uc *HandleWebhookUseCase) handleInstallationDeleted(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	if input.InstallationID <= 0 {
		return nil, domain.ErrInvalidWebhookPayload
	}

	if err := uc.repo.Delete(ctx, input.InstallationID); err != nil {
		return nil, err
	}

	return &HandleWebhookOutput{Message: "installation deleted"}, nil
}

func (uc *HandleWebhookUseCase) handleInstallationSuspend(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	if input.InstallationID <= 0 {
		return nil, domain.ErrInvalidWebhookPayload
	}

	var suspendedAt time.Time
	if input.SuspendedAt != nil {
		parsed, err := time.Parse(time.RFC3339, *input.SuspendedAt)
		if err != nil {
			suspendedAt = time.Now()
		} else {
			suspendedAt = parsed
		}
	} else {
		suspendedAt = time.Now()
	}

	if err := uc.repo.UpdateSuspended(ctx, input.InstallationID, &suspendedAt); err != nil {
		return nil, err
	}

	return &HandleWebhookOutput{Message: "installation suspended"}, nil
}

func (uc *HandleWebhookUseCase) handleInstallationUnsuspend(ctx context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	if input.InstallationID <= 0 {
		return nil, domain.ErrInvalidWebhookPayload
	}

	if err := uc.repo.UpdateSuspended(ctx, input.InstallationID, nil); err != nil {
		return nil, err
	}

	return &HandleWebhookOutput{Message: "installation unsuspended"}, nil
}

func (uc *HandleWebhookUseCase) handleInstallationRepositories(_ context.Context, input HandleWebhookInput) (*HandleWebhookOutput, error) {
	switch input.Action {
	case "added", "removed":
		return &HandleWebhookOutput{Message: "repository change acknowledged"}, nil
	default:
		return &HandleWebhookOutput{Message: "installation_repositories action ignored"}, nil
	}
}
