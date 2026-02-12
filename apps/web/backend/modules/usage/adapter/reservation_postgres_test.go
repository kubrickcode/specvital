package adapter

import (
	"testing"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/domain/port"
)

func TestQuotaReservationPostgresRepository_ImplementsInterface(t *testing.T) {
	var _ port.QuotaReservationRepository = (*QuotaReservationPostgresRepository)(nil)
}

func TestNewQuotaReservationPostgresRepository(t *testing.T) {
	repo := NewQuotaReservationPostgresRepository(nil)
	if repo == nil {
		t.Error("NewQuotaReservationPostgresRepository() returned nil")
	}
}
