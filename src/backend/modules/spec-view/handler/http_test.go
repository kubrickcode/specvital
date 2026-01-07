package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/usecase"
)

type mockConvertSpecViewUseCase struct {
	result *usecase.ConvertSpecViewOutput
	err    error
}

func (m *mockConvertSpecViewUseCase) Execute(_ context.Context, _ usecase.ConvertSpecViewInput) (*usecase.ConvertSpecViewOutput, error) {
	return m.result, m.err
}

func TestConvertSpecView_Success(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	mock := &mockConvertSpecViewUseCase{
		result: &usecase.ConvertSpecViewOutput{
			Files: []usecase.FileConversionResult{
				{
					FilePath:  "src/auth.spec.ts",
					Framework: "vitest",
					Suites: []usecase.SuiteConversionResult{
						{
							Hierarchy: "Auth",
							Name:      "Auth",
							Tests: []usecase.TestConversionResult{
								{
									ConvertedName: "Session creation on login",
									IsFromCache:   true,
									Line:          10,
									OriginalName:  "should create session when login successful",
									Status:        "active",
								},
							},
						},
					},
				},
			},
			Summary: usecase.ConversionSummary{
				CachedCount:    1,
				ConvertedAt:    time.Now(),
				ConvertedCount: 0,
				TotalTests:     1,
			},
		},
	}

	h := &Handler{
		convertSpecView: mock,
		logger:          logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestConvertSpecView_InvalidOwner(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{},
		logger:          logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "invalid owner!",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView400ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestConvertSpecView_EmptyCommitSHA(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{},
		logger:          logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView400ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestConvertSpecView_InvalidCommitSHAFormat(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{},
		logger:          logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "not-a-valid-sha!",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView400ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestConvertSpecView_AnalysisNotFound(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: domain.ErrAnalysisNotFound,
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView404ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 404 response, got %T", resp)
	}
}

func TestConvertSpecView_CommitMismatch(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: domain.ErrCommitMismatch,
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView404ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 404 response, got %T", resp)
	}
}

func TestConvertSpecView_RateLimited(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: domain.ErrRateLimited,
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView429ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 429 response, got %T", resp)
	}
}

func TestConvertSpecView_InvalidLanguage(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: domain.ErrInvalidLanguage,
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView400ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestConvertSpecView_AIProviderUnavailable(t *testing.T) {
	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: domain.ErrAIProviderUnavailable,
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView500ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 500 response, got %T", resp)
	}
}

func TestConvertSpecView_UnexpectedError(t *testing.T) {
	h := &Handler{
		convertSpecView: &mockConvertSpecViewUseCase{
			err: errors.New("unexpected error"),
		},
		logger: logger.New(),
	}

	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView500ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Errorf("expected 500 response, got %T", resp)
	}
}

func TestConvertSpecView_WithLanguage(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	mock := &mockConvertSpecViewUseCase{
		result: &usecase.ConvertSpecViewOutput{
			Files:   []usecase.FileConversionResult{},
			Summary: usecase.ConversionSummary{ConvertedAt: time.Now()},
		},
	}

	h := &Handler{
		convertSpecView: mock,
		logger:          logger.New(),
	}

	lang := api.Ko
	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
		Body: &api.ConvertSpecViewJSONRequestBody{
			Language: &lang,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestConvertSpecView_WithForceRefresh(t *testing.T) {
	t.Skip("Skipped: Spec View feature temporarily suspended pending Level 2 redesign")

	mock := &mockConvertSpecViewUseCase{
		result: &usecase.ConvertSpecViewOutput{
			Files:   []usecase.FileConversionResult{},
			Summary: usecase.ConversionSummary{ConvertedAt: time.Now()},
		},
	}

	h := &Handler{
		convertSpecView: mock,
		logger:          logger.New(),
	}

	forceRefresh := true
	resp, err := h.ConvertSpecView(context.Background(), api.ConvertSpecViewRequestObject{
		Owner:     "test-owner",
		Repo:      "test-repo",
		CommitSHA: "abc1234",
		Body: &api.ConvertSpecViewJSONRequestBody{
			IsForceRefresh: &forceRefresh,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := resp.(api.ConvertSpecView200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestNewHandler_NilConfig(t *testing.T) {
	_, err := NewHandler(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestNewHandler_NilUseCase(t *testing.T) {
	_, err := NewHandler(&HandlerConfig{
		ConvertSpecView: nil,
		Logger:          logger.New(),
	})
	if err == nil {
		t.Error("expected error for nil usecase")
	}
}

func TestNewHandler_NilLogger(t *testing.T) {
	_, err := NewHandler(&HandlerConfig{
		ConvertSpecView: &mockConvertSpecViewUseCase{},
		Logger:          nil,
	})
	if err == nil {
		t.Error("expected error for nil logger")
	}
}

var _ ConvertSpecViewExecutor = (*mockConvertSpecViewUseCase)(nil)
var _ entity.Language = entity.LanguageEn
