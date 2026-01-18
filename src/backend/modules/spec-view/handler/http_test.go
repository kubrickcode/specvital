package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/web/src/backend/common/logger"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type mockTierLookup struct {
	tier string
	err  error
}

func (m *mockTierLookup) GetUserTier(_ context.Context, _ string) (string, error) {
	return m.tier, m.err
}

func TestHandler_lookupUserTier(t *testing.T) {
	log := logger.New()

	t.Run("returns empty tier when userID is empty", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: &mockTierLookup{tier: "pro"},
		}
		got := h.lookupUserTier(context.Background(), "")
		if got != "" {
			t.Errorf("lookupUserTier() = %q, want empty", got)
		}
	})

	t.Run("returns empty tier when tierLookup is nil", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: nil,
		}
		got := h.lookupUserTier(context.Background(), "user-123")
		if got != "" {
			t.Errorf("lookupUserTier() = %q, want empty", got)
		}
	})

	t.Run("returns empty tier on lookup error", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: &mockTierLookup{err: errors.New("db error")},
		}
		got := h.lookupUserTier(context.Background(), "user-123")
		if got != "" {
			t.Errorf("lookupUserTier() = %q, want empty", got)
		}
	})

	t.Run("returns pro tier on successful lookup", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: &mockTierLookup{tier: "pro"},
		}
		got := h.lookupUserTier(context.Background(), "user-123")
		if got != subscription.PlanTierPro {
			t.Errorf("lookupUserTier() = %q, want %q", got, subscription.PlanTierPro)
		}
	})

	t.Run("returns pro_plus tier on successful lookup", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: &mockTierLookup{tier: "pro_plus"},
		}
		got := h.lookupUserTier(context.Background(), "user-123")
		if got != subscription.PlanTierProPlus {
			t.Errorf("lookupUserTier() = %q, want %q", got, subscription.PlanTierProPlus)
		}
	})

	t.Run("returns free tier on successful lookup", func(t *testing.T) {
		h := &Handler{
			logger:     log,
			tierLookup: &mockTierLookup{tier: "free"},
		}
		got := h.lookupUserTier(context.Background(), "user-123")
		if got != subscription.PlanTierFree {
			t.Errorf("lookupUserTier() = %q, want %q", got, subscription.PlanTierFree)
		}
	})
}
