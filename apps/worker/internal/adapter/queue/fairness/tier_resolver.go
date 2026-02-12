package fairness

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kubrickcode/specvital/apps/worker/internal/infra/db"
)

// TierResolver resolves a user's subscription tier.
type TierResolver interface {
	ResolveTier(ctx context.Context, userID string) PlanTier
}

// DBTierResolver resolves user tier from database.
// Returns TierFree if user has no active subscription or on error.
type DBTierResolver struct {
	queries *db.Queries
}

// NewDBTierResolver creates a new DBTierResolver with the given queries.
func NewDBTierResolver(queries *db.Queries) *DBTierResolver {
	return &DBTierResolver{queries: queries}
}

// ResolveTier returns the user's tier from their active subscription.
// Returns TierFree if:
// - userID is empty or invalid UUID
// - user has no active subscription
// - database error occurs
func (r *DBTierResolver) ResolveTier(ctx context.Context, userID string) PlanTier {
	if userID == "" {
		return TierFree
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.WarnContext(ctx, "invalid user ID format, defaulting to free tier",
			"user_id", userID,
			"error", err,
		)
		return TierFree
	}

	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	dbTier, err := r.queries.GetUserTier(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.DebugContext(ctx, "no active subscription found, defaulting to free tier",
				"user_id", userID,
			)
			return TierFree
		}
		slog.WarnContext(ctx, "failed to resolve user tier, defaulting to free tier",
			"user_id", userID,
			"error", err,
		)
		return TierFree
	}

	return convertDBTier(dbTier)
}

// convertDBTier converts db.PlanTier to fairness.PlanTier.
func convertDBTier(dbTier db.PlanTier) PlanTier {
	switch dbTier {
	case db.PlanTierFree:
		return TierFree
	case db.PlanTierPro:
		return TierPro
	case db.PlanTierProPlus:
		return TierProPlus
	case db.PlanTierEnterprise:
		return TierEnterprise
	default:
		return TierFree
	}
}
