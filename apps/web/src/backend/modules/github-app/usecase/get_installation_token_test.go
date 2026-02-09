package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/github-app/domain"
	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

type mockGitHubAppClient struct {
	createTokenFn func(context.Context, int64) (*port.InstallationToken, error)
}

func (m *mockGitHubAppClient) CreateInstallationToken(ctx context.Context, installationID int64) (*port.InstallationToken, error) {
	if m.createTokenFn != nil {
		return m.createTokenFn(ctx, installationID)
	}
	return &port.InstallationToken{
		Token:     "test-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

func (m *mockGitHubAppClient) GetInstallationURL() string {
	return "https://github.com/apps/test"
}

type mockInstallationRepository struct {
	getByInstallationIDFn func(context.Context, int64) (*entity.Installation, error)
}

func (m *mockInstallationRepository) Delete(ctx context.Context, installationID int64) error {
	return nil
}

func (m *mockInstallationRepository) GetByAccountID(ctx context.Context, accountID int64) (*entity.Installation, error) {
	return nil, nil
}

func (m *mockInstallationRepository) GetByInstallationID(ctx context.Context, installationID int64) (*entity.Installation, error) {
	if m.getByInstallationIDFn != nil {
		return m.getByInstallationIDFn(ctx, installationID)
	}
	return nil, domain.ErrInstallationNotFound
}

func (m *mockInstallationRepository) ListByAccountIDs(ctx context.Context, accountIDs []int64) ([]entity.Installation, error) {
	return nil, nil
}

func (m *mockInstallationRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Installation, error) {
	return nil, nil
}

func (m *mockInstallationRepository) ListOrganizations(ctx context.Context) ([]entity.Installation, error) {
	return nil, nil
}

func (m *mockInstallationRepository) UpdateSuspended(ctx context.Context, installationID int64, suspendedAt *time.Time) error {
	return nil
}

func (m *mockInstallationRepository) Upsert(ctx context.Context, installation *entity.Installation) error {
	return nil
}

func TestGetInstallationTokenUseCase_Execute_Success(t *testing.T) {
	repo := &mockInstallationRepository{
		getByInstallationIDFn: func(ctx context.Context, id int64) (*entity.Installation, error) {
			return &entity.Installation{
				InstallationID: id,
				AccountType:    entity.AccountTypeOrganization,
				AccountLogin:   "test-org",
			}, nil
		},
	}

	ghClient := &mockGitHubAppClient{
		createTokenFn: func(ctx context.Context, id int64) (*port.InstallationToken, error) {
			return &port.InstallationToken{
				Token:     "generated-token",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			}, nil
		},
	}

	uc := NewGetInstallationTokenUseCase(ghClient, repo)

	output, err := uc.Execute(context.Background(), GetInstallationTokenInput{
		InstallationID: 123,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Token != "generated-token" {
		t.Errorf("expected token 'generated-token', got '%s'", output.Token)
	}
}

func TestGetInstallationTokenUseCase_Execute_InstallationNotFound(t *testing.T) {
	repo := &mockInstallationRepository{
		getByInstallationIDFn: func(ctx context.Context, id int64) (*entity.Installation, error) {
			return nil, domain.ErrInstallationNotFound
		},
	}

	ghClient := &mockGitHubAppClient{}

	uc := NewGetInstallationTokenUseCase(ghClient, repo)

	_, err := uc.Execute(context.Background(), GetInstallationTokenInput{
		InstallationID: 123,
	})

	if !errors.Is(err, domain.ErrInstallationNotFound) {
		t.Errorf("expected ErrInstallationNotFound, got %v", err)
	}
}

func TestGetInstallationTokenUseCase_Execute_TokenGenerationFailed(t *testing.T) {
	repo := &mockInstallationRepository{
		getByInstallationIDFn: func(ctx context.Context, id int64) (*entity.Installation, error) {
			return &entity.Installation{
				InstallationID: id,
				AccountType:    entity.AccountTypeOrganization,
			}, nil
		},
	}

	ghClient := &mockGitHubAppClient{
		createTokenFn: func(ctx context.Context, id int64) (*port.InstallationToken, error) {
			return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, errors.New("API error"))
		},
	}

	uc := NewGetInstallationTokenUseCase(ghClient, repo)

	_, err := uc.Execute(context.Background(), GetInstallationTokenInput{
		InstallationID: 123,
	})

	if !errors.Is(err, domain.ErrTokenGenerationFailed) {
		t.Errorf("expected ErrTokenGenerationFailed, got %v", err)
	}
}

func TestGetInstallationTokenUseCase_Execute_InstallationSuspended(t *testing.T) {
	suspendedAt := time.Now()
	repo := &mockInstallationRepository{
		getByInstallationIDFn: func(ctx context.Context, id int64) (*entity.Installation, error) {
			return &entity.Installation{
				InstallationID: id,
				AccountType:    entity.AccountTypeOrganization,
				SuspendedAt:    &suspendedAt,
			}, nil
		},
	}

	ghClient := &mockGitHubAppClient{}

	uc := NewGetInstallationTokenUseCase(ghClient, repo)

	_, err := uc.Execute(context.Background(), GetInstallationTokenInput{
		InstallationID: 123,
	})

	if !errors.Is(err, domain.ErrInstallationSuspended) {
		t.Errorf("expected ErrInstallationSuspended, got %v", err)
	}
}
