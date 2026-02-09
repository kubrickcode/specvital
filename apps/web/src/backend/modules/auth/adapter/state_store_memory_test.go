package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/auth/domain"
)

func TestMemoryStateStore_CreateAndValidate(t *testing.T) {
	store := NewMemoryStateStore()
	ctx := context.Background()

	state, err := store.Create(ctx)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if state == "" {
		t.Fatal("Create() returned empty state")
	}

	if len(state) < 32 {
		t.Errorf("Create() state too short, got %d characters", len(state))
	}

	if err := store.Validate(ctx, state); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestMemoryStateStore_ValidateConsumed(t *testing.T) {
	store := NewMemoryStateStore()
	ctx := context.Background()

	state, err := store.Create(ctx)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := store.Validate(ctx, state); err != nil {
		t.Fatalf("first Validate() error = %v", err)
	}

	if err := store.Validate(ctx, state); err != domain.ErrInvalidState {
		t.Errorf("second Validate() expected ErrInvalidState, got %v", err)
	}
}

func TestMemoryStateStore_ValidateInvalid(t *testing.T) {
	store := NewMemoryStateStore()
	ctx := context.Background()

	if err := store.Validate(ctx, "invalid-state"); err != domain.ErrInvalidState {
		t.Errorf("Validate() expected ErrInvalidState, got %v", err)
	}
}

func TestMemoryStateStore_ValidateExpired(t *testing.T) {
	store := NewMemoryStateStoreWithTTL(10 * time.Millisecond)
	ctx := context.Background()

	state, err := store.Create(ctx)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	if err := store.Validate(ctx, state); err != domain.ErrInvalidState {
		t.Errorf("Validate() expected ErrInvalidState for expired state, got %v", err)
	}
}

func TestMemoryStateStore_UniqueStates(t *testing.T) {
	store := NewMemoryStateStore()
	ctx := context.Background()

	states := make(map[string]bool)
	for i := 0; i < 100; i++ {
		state, err := store.Create(ctx)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if states[state] {
			t.Fatalf("Create() produced duplicate state: %s", state)
		}
		states[state] = true
	}
}
