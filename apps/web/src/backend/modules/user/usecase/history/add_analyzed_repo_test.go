package history

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type mockAddHistoryRepository struct {
	result *port.AddHistoryResult
	err    error
}

func (m *mockAddHistoryRepository) AddUserAnalyzedRepository(_ context.Context, _, _, _ string) (*port.AddHistoryResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockAddHistoryRepository) CheckUserHistoryExists(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockAddHistoryRepository) GetUserAnalyzedRepositories(_ context.Context, _ string, _ port.AnalyzedReposParams) ([]*entity.AnalyzedRepository, error) {
	return nil, nil
}

func TestAddAnalyzedRepoUseCase_Success(t *testing.T) {
	now := time.Now()
	repo := &mockAddHistoryRepository{
		result: &port.AddHistoryResult{
			AnalysisID: "analysis-123",
			UpdatedAt:  now,
		},
	}
	uc := NewAddAnalyzedRepoUseCase(repo)

	output, err := uc.Execute(context.Background(), AddAnalyzedRepoInput{
		Owner:  "facebook",
		Repo:   "react",
		UserID: "user-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.AnalysisID != "analysis-123" {
		t.Errorf("expected AnalysisID 'analysis-123', got %s", output.AnalysisID)
	}
	if !output.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt %v, got %v", now, output.UpdatedAt)
	}
}

func TestAddAnalyzedRepoUseCase_NotFound(t *testing.T) {
	repo := &mockAddHistoryRepository{
		err: domain.ErrAnalysisNotFound,
	}
	uc := NewAddAnalyzedRepoUseCase(repo)

	_, err := uc.Execute(context.Background(), AddAnalyzedRepoInput{
		Owner:  "unknown",
		Repo:   "repo",
		UserID: "user-123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAddAnalyzedRepoUseCase_RepositoryError(t *testing.T) {
	repo := &mockAddHistoryRepository{
		err: errors.New("database connection failed"),
	}
	uc := NewAddAnalyzedRepoUseCase(repo)

	_, err := uc.Execute(context.Background(), AddAnalyzedRepoInput{
		Owner:  "facebook",
		Repo:   "react",
		UserID: "user-123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
