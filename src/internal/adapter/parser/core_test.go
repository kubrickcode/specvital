package parser

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/specvital/worker/internal/domain/analysis"
)

func TestNewCoreParser(t *testing.T) {
	parser := NewCoreParser()
	if parser == nil {
		t.Fatal("NewCoreParser returned nil")
	}
}

func TestCoreParser_Scan_InvalidSourceType(t *testing.T) {
	p := NewCoreParser()

	// mockSource doesn't implement coreSourceProvider
	mockSrc := &mockInvalidSource{}

	_, err := p.Scan(context.Background(), mockSrc)
	if err == nil {
		t.Fatal("expected error for source not implementing coreSourceProvider")
	}
	if !strings.Contains(err.Error(), "does not implement coreSourceProvider interface") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCoreParser_ScanStream_InvalidSourceType(t *testing.T) {
	p := NewCoreParser()

	// mockSource doesn't implement coreSourceProvider
	mockSrc := &mockInvalidSource{}

	_, err := p.ScanStream(context.Background(), mockSrc)
	if err == nil {
		t.Fatal("expected error for source not implementing coreSourceProvider")
	}
	if !strings.Contains(err.Error(), "does not implement coreSourceProvider interface") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// mockInvalidSource implements analysis.Source but not coreSourceProvider.
type mockInvalidSource struct{}

var _ analysis.Source = (*mockInvalidSource)(nil)

func (m *mockInvalidSource) Branch() string                { return "" }
func (m *mockInvalidSource) CommitSHA() string             { return "" }
func (m *mockInvalidSource) CommittedAt() time.Time        { return time.Time{} }
func (m *mockInvalidSource) Close(_ context.Context) error { return nil }

func (m *mockInvalidSource) VerifyCommitExists(_ context.Context, _ string) (bool, error) {
	return true, nil
}

// Conversion tests moved to adapter/mapping/core_domain_test.go
