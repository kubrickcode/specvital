package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/github-app/domain"
	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

var _ port.InstallationRepository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) Delete(ctx context.Context, installationID int64) error {
	return r.queries.DeleteGitHubAppInstallation(ctx, installationID)
}

func (r *PostgresRepository) GetByAccountID(ctx context.Context, accountID int64) (*entity.Installation, error) {
	if accountID <= 0 {
		return nil, domain.ErrInstallationNotFound
	}

	row, err := r.queries.GetGitHubAppInstallationByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInstallationNotFound
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *PostgresRepository) GetByInstallationID(ctx context.Context, installationID int64) (*entity.Installation, error) {
	if installationID <= 0 {
		return nil, domain.ErrInstallationNotFound
	}

	row, err := r.queries.GetGitHubAppInstallationByID(ctx, installationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInstallationNotFound
		}
		return nil, err
	}
	return toEntity(row), nil
}

func (r *PostgresRepository) ListByAccountIDs(ctx context.Context, accountIDs []int64) ([]entity.Installation, error) {
	if len(accountIDs) == 0 {
		return []entity.Installation{}, nil
	}

	rows, err := r.queries.ListGitHubAppInstallationsByAccountIDs(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	installations := make([]entity.Installation, 0, len(rows))
	for _, row := range rows {
		installations = append(installations, *toEntity(row))
	}
	return installations, nil
}

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Installation, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.ListGitHubAppInstallationsByUserID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return nil, err
	}

	installations := make([]entity.Installation, 0, len(rows))
	for _, row := range rows {
		installations = append(installations, *toEntity(row))
	}
	return installations, nil
}

func (r *PostgresRepository) ListOrganizations(ctx context.Context) ([]entity.Installation, error) {
	rows, err := r.queries.ListGitHubAppOrganizationInstallations(ctx)
	if err != nil {
		return nil, err
	}

	installations := make([]entity.Installation, 0, len(rows))
	for _, row := range rows {
		installations = append(installations, *toEntity(row))
	}
	return installations, nil
}

func (r *PostgresRepository) UpdateSuspended(ctx context.Context, installationID int64, suspendedAt *time.Time) error {
	var pgSuspendedAt pgtype.Timestamptz
	if suspendedAt != nil {
		pgSuspendedAt = pgtype.Timestamptz{Time: *suspendedAt, Valid: true}
	}

	return r.queries.UpdateGitHubAppInstallationSuspended(ctx, db.UpdateGitHubAppInstallationSuspendedParams{
		InstallationID: installationID,
		SuspendedAt:    pgSuspendedAt,
	})
}

func (r *PostgresRepository) Upsert(ctx context.Context, installation *entity.Installation) error {
	var installerUserID pgtype.UUID
	if installation.InstallerUserID != nil {
		uid, err := uuid.Parse(*installation.InstallerUserID)
		if err != nil {
			return err
		}
		installerUserID = pgtype.UUID{Bytes: uid, Valid: true}
	}

	var accountAvatarURL pgtype.Text
	if installation.AccountAvatarURL != nil {
		accountAvatarURL = pgtype.Text{String: *installation.AccountAvatarURL, Valid: true}
	}

	var suspendedAt pgtype.Timestamptz
	if installation.SuspendedAt != nil {
		suspendedAt = pgtype.Timestamptz{Time: *installation.SuspendedAt, Valid: true}
	}

	return r.queries.UpsertGitHubAppInstallation(ctx, db.UpsertGitHubAppInstallationParams{
		InstallationID:   installation.InstallationID,
		AccountType:      db.GithubAccountType(installation.AccountType),
		AccountID:        installation.AccountID,
		AccountLogin:     installation.AccountLogin,
		AccountAvatarUrl: accountAvatarURL,
		InstallerUserID:  installerUserID,
		SuspendedAt:      suspendedAt,
	})
}

func toEntity(row db.GithubAppInstallation) *entity.Installation {
	var accountAvatarURL *string
	if row.AccountAvatarUrl.Valid {
		accountAvatarURL = &row.AccountAvatarUrl.String
	}

	var installerUserID *string
	if row.InstallerUserID.Valid {
		uid := uuid.UUID(row.InstallerUserID.Bytes).String()
		installerUserID = &uid
	}

	var suspendedAt *time.Time
	if row.SuspendedAt.Valid {
		suspendedAt = &row.SuspendedAt.Time
	}

	return &entity.Installation{
		AccountAvatarURL: accountAvatarURL,
		AccountID:        row.AccountID,
		AccountLogin:     row.AccountLogin,
		AccountType:      entity.AccountType(row.AccountType),
		CreatedAt:        row.CreatedAt.Time,
		ID:               uuid.UUID(row.ID.Bytes).String(),
		InstallationID:   row.InstallationID,
		InstallerUserID:  installerUserID,
		SuspendedAt:      suspendedAt,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}
