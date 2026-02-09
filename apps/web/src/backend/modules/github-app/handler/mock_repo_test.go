package handler

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/github-app/domain"
	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

var _ port.InstallationRepository = (*mockRepo)(nil)

type mockRepo struct {
	installations map[int64]*entity.Installation
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		installations: make(map[int64]*entity.Installation),
	}
}

func (r *mockRepo) Delete(_ context.Context, installationID int64) error {
	delete(r.installations, installationID)
	return nil
}

func (r *mockRepo) GetByAccountID(_ context.Context, accountID int64) (*entity.Installation, error) {
	for _, inst := range r.installations {
		if inst.AccountID == accountID {
			return inst, nil
		}
	}
	return nil, domain.ErrInstallationNotFound
}

func (r *mockRepo) GetByInstallationID(_ context.Context, installationID int64) (*entity.Installation, error) {
	if inst, ok := r.installations[installationID]; ok {
		return inst, nil
	}
	return nil, domain.ErrInstallationNotFound
}

func (r *mockRepo) ListByAccountIDs(_ context.Context, accountIDs []int64) ([]entity.Installation, error) {
	var result []entity.Installation
	for _, accountID := range accountIDs {
		for _, inst := range r.installations {
			if inst.AccountID == accountID {
				result = append(result, *inst)
				break
			}
		}
	}
	return result, nil
}

func (r *mockRepo) ListByUserID(_ context.Context, _ string) ([]entity.Installation, error) {
	return nil, nil
}

func (r *mockRepo) ListOrganizations(_ context.Context) ([]entity.Installation, error) {
	var result []entity.Installation
	for _, inst := range r.installations {
		if inst.AccountType == entity.AccountTypeOrganization {
			result = append(result, *inst)
		}
	}
	return result, nil
}

func (r *mockRepo) UpdateSuspended(_ context.Context, installationID int64, suspendedAt *time.Time) error {
	if inst, ok := r.installations[installationID]; ok {
		inst.SuspendedAt = suspendedAt
	}
	return nil
}

func (r *mockRepo) Upsert(_ context.Context, installation *entity.Installation) error {
	r.installations[installation.InstallationID] = installation
	return nil
}
