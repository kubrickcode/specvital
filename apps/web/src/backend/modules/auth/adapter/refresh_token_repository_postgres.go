package adapter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type RefreshTokenPostgresRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var _ port.RefreshTokenRepository = (*RefreshTokenPostgresRepository)(nil)

func NewRefreshTokenPostgresRepository(pool *pgxpool.Pool, queries *db.Queries) *RefreshTokenPostgresRepository {
	if pool == nil {
		panic("pool is required")
	}
	if queries == nil {
		panic("queries is required")
	}
	return &RefreshTokenPostgresRepository{pool: pool, queries: queries}
}

func (r *RefreshTokenPostgresRepository) Create(ctx context.Context, token *entity.RefreshToken) (string, error) {
	if token == nil {
		return "", errors.New("token is required")
	}
	if token.TokenHash == "" {
		return "", errors.New("token hash is required")
	}
	if token.UserID == "" {
		return "", errors.New("user ID is required")
	}
	if token.FamilyID == "" {
		return "", errors.New("family ID is required")
	}

	userUUID, err := stringToUUID(token.UserID)
	if err != nil {
		return "", fmt.Errorf("parse user ID: %w", err)
	}

	familyUUID, err := stringToUUID(token.FamilyID)
	if err != nil {
		return "", fmt.Errorf("parse family ID: %w", err)
	}

	params := db.CreateRefreshTokenParams{
		UserID:    userUUID,
		TokenHash: token.TokenHash,
		FamilyID:  familyUUID,
		ExpiresAt: pgtype.Timestamptz{Time: token.ExpiresAt, Valid: true},
	}

	if token.Replaces != nil {
		replacesUUID, err := stringToUUID(*token.Replaces)
		if err != nil {
			return "", fmt.Errorf("parse replaces ID: %w", err)
		}
		params.Replaces = replacesUUID
	}

	row, err := r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return "", fmt.Errorf("create refresh token: %w", err)
	}

	return uuidToString(row.ID), nil
}

func (r *RefreshTokenPostgresRepository) GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	row, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}

	return mapRefreshTokenFromDB(&row), nil
}

func (r *RefreshTokenPostgresRepository) Revoke(ctx context.Context, id string) error {
	uuid, err := stringToUUID(id)
	if err != nil {
		return fmt.Errorf("parse token ID: %w", err)
	}

	rowsAffected, err := r.queries.RevokeRefreshToken(ctx, uuid)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenPostgresRepository) RevokeFamily(ctx context.Context, familyID string) error {
	uuid, err := stringToUUID(familyID)
	if err != nil {
		return fmt.Errorf("parse family ID: %w", err)
	}

	if err := r.queries.RevokeRefreshTokenFamily(ctx, uuid); err != nil {
		return fmt.Errorf("revoke refresh token family: %w", err)
	}

	return nil
}

func (r *RefreshTokenPostgresRepository) RevokeUserTokens(ctx context.Context, userID string) error {
	uuid, err := stringToUUID(userID)
	if err != nil {
		return fmt.Errorf("parse user ID: %w", err)
	}

	if err := r.queries.RevokeUserRefreshTokens(ctx, uuid); err != nil {
		return fmt.Errorf("revoke user refresh tokens: %w", err)
	}

	return nil
}

func (r *RefreshTokenPostgresRepository) RotateToken(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error) {
	if oldTokenID == "" {
		return "", errors.New("old token ID is required")
	}
	if newToken == nil {
		return "", errors.New("new token is required")
	}

	var newTokenID string
	err := r.withTx(ctx, func(qtx *db.Queries) error {
		oldUUID, err := stringToUUID(oldTokenID)
		if err != nil {
			return fmt.Errorf("parse old token ID: %w", err)
		}

		rowsAffected, err := qtx.RevokeRefreshToken(ctx, oldUUID)
		if err != nil {
			return fmt.Errorf("revoke old token: %w", err)
		}
		if rowsAffected == 0 {
			return domain.ErrRefreshTokenNotFound
		}

		userUUID, err := stringToUUID(newToken.UserID)
		if err != nil {
			return fmt.Errorf("parse user ID: %w", err)
		}
		familyUUID, err := stringToUUID(newToken.FamilyID)
		if err != nil {
			return fmt.Errorf("parse family ID: %w", err)
		}

		params := db.CreateRefreshTokenParams{
			UserID:    userUUID,
			TokenHash: newToken.TokenHash,
			FamilyID:  familyUUID,
			ExpiresAt: pgtype.Timestamptz{Time: newToken.ExpiresAt, Valid: true},
		}

		if newToken.Replaces != nil {
			replacesUUID, err := stringToUUID(*newToken.Replaces)
			if err != nil {
				return fmt.Errorf("parse replaces ID: %w", err)
			}
			params.Replaces = replacesUUID
		}

		row, err := qtx.CreateRefreshToken(ctx, params)
		if err != nil {
			return fmt.Errorf("create new token: %w", err)
		}

		newTokenID = uuidToString(row.ID)
		return nil
	})
	if err != nil {
		return "", err
	}

	return newTokenID, nil
}

func (r *RefreshTokenPostgresRepository) withTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(r.queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func mapRefreshTokenFromDB(row *db.RefreshToken) *entity.RefreshToken {
	token := &entity.RefreshToken{
		CreatedAt: row.CreatedAt.Time,
		ExpiresAt: row.ExpiresAt.Time,
		FamilyID:  uuidToString(row.FamilyID),
		ID:        uuidToString(row.ID),
		TokenHash: row.TokenHash,
		UserID:    uuidToString(row.UserID),
	}

	if row.RevokedAt.Valid {
		t := row.RevokedAt.Time
		token.RevokedAt = &t
	}
	if row.Replaces.Valid {
		s := uuidToString(row.Replaces)
		token.Replaces = &s
	}

	return token
}
