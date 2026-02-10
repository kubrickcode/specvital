package handler

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/apps/web/src/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/domain/port"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/usecase"
)

func TestNewAPIHandler_ValidConfig(t *testing.T) {
	mockClient := &mockGitHubAppClient{installURL: "https://example.com"}
	mockRepo := newMockRepo()

	getInstallURLUC := usecase.NewGetInstallURLUseCase(mockClient)
	listInstallationsUC := usecase.NewListInstallationsUseCase(mockRepo)

	h, err := NewAPIHandler(&APIHandlerConfig{
		GetInstallURL:     getInstallURLUC,
		ListInstallations: listInstallationsUC,
		Logger:            logger.New(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Error("expected handler to be non-nil")
	}
}

func TestNewAPIHandler_NilConfig(t *testing.T) {
	_, err := NewAPIHandler(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestNewAPIHandler_MissingGetInstallURL(t *testing.T) {
	mockRepo := newMockRepo()
	listInstallationsUC := usecase.NewListInstallationsUseCase(mockRepo)

	_, err := NewAPIHandler(&APIHandlerConfig{
		GetInstallURL:     nil,
		ListInstallations: listInstallationsUC,
		Logger:            logger.New(),
	})
	if err == nil {
		t.Error("expected error for missing GetInstallURL")
	}
}

func TestNewAPIHandler_MissingListInstallations(t *testing.T) {
	mockClient := &mockGitHubAppClient{installURL: "https://example.com"}
	getInstallURLUC := usecase.NewGetInstallURLUseCase(mockClient)

	_, err := NewAPIHandler(&APIHandlerConfig{
		GetInstallURL:     getInstallURLUC,
		ListInstallations: nil,
		Logger:            logger.New(),
	})
	if err == nil {
		t.Error("expected error for missing ListInstallations")
	}
}

func TestNewAPIHandler_MissingLogger(t *testing.T) {
	mockClient := &mockGitHubAppClient{installURL: "https://example.com"}
	mockRepo := newMockRepo()

	getInstallURLUC := usecase.NewGetInstallURLUseCase(mockClient)
	listInstallationsUC := usecase.NewListInstallationsUseCase(mockRepo)

	_, err := NewAPIHandler(&APIHandlerConfig{
		GetInstallURL:     getInstallURLUC,
		ListInstallations: listInstallationsUC,
		Logger:            nil,
	})
	if err == nil {
		t.Error("expected error for missing logger")
	}
}

type mockGitHubAppClient struct {
	installURL string
}

var _ port.GitHubAppClient = (*mockGitHubAppClient)(nil)

func (m *mockGitHubAppClient) CreateInstallationToken(_ context.Context, _ int64) (*port.InstallationToken, error) {
	return nil, nil
}

func (m *mockGitHubAppClient) GetInstallationURL() string {
	return m.installURL
}
