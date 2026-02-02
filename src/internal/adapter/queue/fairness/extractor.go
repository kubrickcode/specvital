package fairness

import (
	"encoding/json"
	"log/slog"
)

const maxArgsByteSize = 64 * 1024 // 64KB - reasonable limit for job args

// UserJobExtractor extracts user identity from River job args.
type UserJobExtractor interface {
	ExtractUserID(encodedArgs []byte) string
}

// JSONArgsExtractor extracts user_id and tier from JSON-encoded job arguments.
// Compatible with both SpecView and Analysis worker args.
type JSONArgsExtractor struct{}

// NewJSONArgsExtractor creates a new JSONArgsExtractor.
func NewJSONArgsExtractor() *JSONArgsExtractor {
	return &JSONArgsExtractor{}
}

// ExtractUserID parses the user_id field from JSON-encoded args.
// Returns empty string if user_id is missing, invalid JSON, or args exceed size limit.
func (e *JSONArgsExtractor) ExtractUserID(encodedArgs []byte) string {
	if len(encodedArgs) > maxArgsByteSize {
		slog.Warn("job args exceed size limit, skipping fairness",
			"args_size", len(encodedArgs),
			"max_size", maxArgsByteSize,
		)
		return ""
	}

	var args struct {
		UserID *string `json:"user_id"` // pointer to detect null vs missing
	}

	if err := json.Unmarshal(encodedArgs, &args); err != nil {
		slog.Warn("failed to parse user_id from job args, skipping fairness",
			"error", err,
			"args_size", len(encodedArgs),
		)
		return ""
	}

	if args.UserID == nil {
		return ""
	}

	return *args.UserID
}

// ExtractTier parses the tier field from JSON-encoded args.
// Defaults to TierFree if tier is missing, invalid, unknown, or args exceed size limit.
//
// Deprecated: Use TierResolver.ResolveTier instead for DB-based tier lookup.
// This method is kept for backward compatibility and will be removed in a future version.
func (e *JSONArgsExtractor) ExtractTier(encodedArgs []byte) PlanTier {
	if len(encodedArgs) > maxArgsByteSize {
		slog.Warn("job args exceed size limit, defaulting to free tier",
			"args_size", len(encodedArgs),
			"max_size", maxArgsByteSize,
		)
		return TierFree
	}

	var args struct {
		Tier string `json:"tier"`
	}

	if err := json.Unmarshal(encodedArgs, &args); err != nil {
		slog.Warn("failed to parse tier from job args, defaulting to free",
			"error", err,
			"args_size", len(encodedArgs),
		)
		return TierFree
	}

	tier := PlanTier(args.Tier)

	// Validate against known tiers
	switch tier {
	case TierFree, TierPro, TierProPlus, TierEnterprise:
		return tier
	default:
		if args.Tier != "" {
			slog.Warn("unknown tier, defaulting to free",
				"tier", args.Tier,
			)
		}
		return TierFree
	}
}
