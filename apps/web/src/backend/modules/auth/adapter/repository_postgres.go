package adapter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type PostgresRepository struct {
	queries *db.Queries
}

var _ port.Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	if queries == nil {
		panic("queries is required")
	}
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	params := db.CreateUserParams{
		Username: user.Username,
	}

	if user.Email != nil {
		params.Email = pgtype.Text{String: *user.Email, Valid: true}
	}
	if user.AvatarURL != "" {
		params.AvatarUrl = pgtype.Text{String: user.AvatarURL, Valid: true}
	}

	id, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	return uuidToString(id), nil
}

func (r *PostgresRepository) GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error) {
	params := db.GetOAuthAccountByProviderParams{
		Provider:       db.OauthProvider(provider),
		ProviderUserID: providerUserID,
	}

	row, err := r.queries.GetOAuthAccountByProvider(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get oauth account by provider: %w", err)
	}

	return mapOAuthAccountFromDB(&row), nil
}

func (r *PostgresRepository) GetOAuthAccountByUserID(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error) {
	uuid, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	params := db.GetOAuthAccountByUserIDParams{
		UserID:   uuid,
		Provider: db.OauthProvider(provider),
	}

	row, err := r.queries.GetOAuthAccountByUserID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get oauth account by user ID: %w", err)
	}

	return mapOAuthAccountFromDB(&row), nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	uuid, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	row, err := r.queries.GetUserByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return mapUserFromDB(&row), nil
}

func (r *PostgresRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	uuid, err := stringToUUID(userID)
	if err != nil {
		return fmt.Errorf("parse user ID: %w", err)
	}

	if err := r.queries.UpdateLastLogin(ctx, uuid); err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	return nil
}

func (r *PostgresRepository) UpsertOAuthAccount(ctx context.Context, account *entity.OAuthAccount) (string, error) {
	userUUID, err := stringToUUID(account.UserID)
	if err != nil {
		return "", fmt.Errorf("parse user ID: %w", err)
	}

	params := db.UpsertOAuthAccountParams{
		UserID:         userUUID,
		Provider:       db.OauthProvider(account.Provider),
		ProviderUserID: account.ProviderUserID,
	}

	if account.ProviderUsername != nil {
		params.ProviderUsername = pgtype.Text{String: *account.ProviderUsername, Valid: true}
	}
	if account.AccessToken != nil {
		params.AccessToken = pgtype.Text{String: *account.AccessToken, Valid: true}
	}
	if account.Scope != nil {
		params.Scope = pgtype.Text{String: *account.Scope, Valid: true}
	}

	id, err := r.queries.UpsertOAuthAccount(ctx, params)
	if err != nil {
		return "", fmt.Errorf("upsert oauth account: %w", err)
	}

	return uuidToString(id), nil
}

func mapOAuthAccountFromDB(row *db.OauthAccount) *entity.OAuthAccount {
	account := &entity.OAuthAccount{
		CreatedAt:      row.CreatedAt.Time,
		ID:             uuidToString(row.ID),
		Provider:       string(row.Provider),
		ProviderUserID: row.ProviderUserID,
		UpdatedAt:      row.UpdatedAt.Time,
		UserID:         uuidToString(row.UserID),
	}

	if row.AccessToken.Valid {
		account.AccessToken = &row.AccessToken.String
	}
	if row.ProviderUsername.Valid {
		account.ProviderUsername = &row.ProviderUsername.String
	}
	if row.Scope.Valid {
		account.Scope = &row.Scope.String
	}

	return account
}

func mapUserFromDB(row *db.GetUserByIDRow) *entity.User {
	user := &entity.User{
		CreatedAt: row.CreatedAt.Time,
		ID:        uuidToString(row.ID),
		UpdatedAt: row.UpdatedAt.Time,
		Username:  row.Username,
	}

	if row.AvatarUrl.Valid {
		user.AvatarURL = row.AvatarUrl.String
	}
	if row.Email.Valid {
		user.Email = &row.Email.String
	}
	if row.LastLoginAt.Valid {
		t := row.LastLoginAt.Time
		user.LastLoginAt = &t
	}

	return user
}

func stringToUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(s); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
