package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pkgerrors "github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	subscriptionDomain "github.com/specvital/web/src/backend/modules/subscription/domain"
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

func TestErrNoActiveSubscription_WrappedErrorMatching(t *testing.T) {
	t.Run("matches wrapped ErrNoActiveSubscription with fmt.Errorf", func(t *testing.T) {
		wrappedErr := fmt.Errorf("check quota for user xyz: %w", subscriptionDomain.ErrNoActiveSubscription)
		if !pkgerrors.Is(wrappedErr, subscriptionDomain.ErrNoActiveSubscription) {
			t.Error("errors.Is should match wrapped ErrNoActiveSubscription")
		}
	})

	t.Run("matches deeply wrapped ErrNoActiveSubscription", func(t *testing.T) {
		inner := fmt.Errorf("inner: %w", subscriptionDomain.ErrNoActiveSubscription)
		outer := fmt.Errorf("outer: %w", inner)
		if !pkgerrors.Is(outer, subscriptionDomain.ErrNoActiveSubscription) {
			t.Error("errors.Is should match deeply wrapped ErrNoActiveSubscription")
		}
	})
}
