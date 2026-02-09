package history

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type mockAnalysisHistoryRepository struct {
	repos []*entity.AnalyzedRepository
	err   error
}

func (m *mockAnalysisHistoryRepository) AddUserAnalyzedRepository(_ context.Context, _, _, _ string) (*port.AddHistoryResult, error) {
	return nil, nil
}

func (m *mockAnalysisHistoryRepository) CheckUserHistoryExists(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockAnalysisHistoryRepository) GetUserAnalyzedRepositories(_ context.Context, _ string, _ port.AnalyzedReposParams) ([]*entity.AnalyzedRepository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.repos, nil
}

func TestGetAnalyzedReposUseCase_Empty(t *testing.T) {
	repo := &mockAnalysisHistoryRepository{repos: []*entity.AnalyzedRepository{}}
	uc := NewGetAnalyzedReposUseCase(repo)

	result, err := uc.Execute(context.Background(), GetAnalyzedReposInput{
		Limit:     20,
		Ownership: entity.OwnershipAll,
		UserID:    "user-id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasNext {
		t.Error("expected HasNext to be false")
	}
	if result.NextCursor != nil {
		t.Error("expected NextCursor to be nil")
	}
	if len(result.Data) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Data))
	}
}

func TestGetAnalyzedReposUseCase_WithPagination(t *testing.T) {
	now := time.Now()
	repos := make([]*entity.AnalyzedRepository, 21)
	for i := 0; i < 21; i++ {
		repos[i] = &entity.AnalyzedRepository{
			CodebaseID:  "codebase-" + strconv.Itoa(i),
			CommitSHA:   "sha-" + strconv.Itoa(i),
			CompletedAt: now,
			HistoryID:   "history-" + strconv.Itoa(i),
			Name:        "repo",
			Owner:       "owner",
			TotalTests:  100,
			UpdatedAt:   now.Add(-time.Duration(i) * time.Hour),
		}
	}

	repo := &mockAnalysisHistoryRepository{repos: repos}
	uc := NewGetAnalyzedReposUseCase(repo)

	result, err := uc.Execute(context.Background(), GetAnalyzedReposInput{
		Limit:     20,
		Ownership: entity.OwnershipAll,
		UserID:    "user-id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasNext {
		t.Error("expected HasNext to be true")
	}
	if result.NextCursor == nil {
		t.Error("expected NextCursor to be set")
	}
	if len(result.Data) != 20 {
		t.Errorf("expected 20 items, got %d", len(result.Data))
	}
}

func TestGetAnalyzedReposUseCase_ExactLimit(t *testing.T) {
	now := time.Now()
	repos := make([]*entity.AnalyzedRepository, 20)
	for i := 0; i < 20; i++ {
		repos[i] = &entity.AnalyzedRepository{
			CodebaseID:  "codebase-" + strconv.Itoa(i),
			CommitSHA:   "sha-" + strconv.Itoa(i),
			CompletedAt: now,
			HistoryID:   "history-" + strconv.Itoa(i),
			Name:        "repo",
			Owner:       "owner",
			TotalTests:  100,
			UpdatedAt:   now.Add(-time.Duration(i) * time.Hour),
		}
	}

	repo := &mockAnalysisHistoryRepository{repos: repos}
	uc := NewGetAnalyzedReposUseCase(repo)

	result, err := uc.Execute(context.Background(), GetAnalyzedReposInput{
		Limit:     20,
		Ownership: entity.OwnershipAll,
		UserID:    "user-id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasNext {
		t.Error("expected HasNext to be false")
	}
	if result.NextCursor != nil {
		t.Error("expected NextCursor to be nil")
	}
	if len(result.Data) != 20 {
		t.Errorf("expected 20 items, got %d", len(result.Data))
	}
}

func TestParseOwnershipFilter(t *testing.T) {
	tests := []struct {
		input    string
		expected entity.OwnershipFilter
	}{
		{"all", entity.OwnershipAll},
		{"mine", entity.OwnershipMine},
		{"organization", entity.OwnershipOrganization},
		{"", entity.OwnershipAll},
		{"invalid", entity.OwnershipAll},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := entity.ParseOwnershipFilter(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
