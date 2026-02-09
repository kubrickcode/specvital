package handler

import (
	"context"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/github/domain"
	ghentity "github.com/specvital/web/src/backend/modules/github/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/usecase"
)

func contextWithUserID(userID string) context.Context {
	claims := &entity.Claims{
		ExpiresAt: time.Now().Add(time.Hour),
		IssuedAt:  time.Now(),
		Issuer:    "specvital",
		Subject:   userID,
	}
	return middleware.WithClaims(context.Background(), claims)
}

type mockListUserReposUseCase struct {
	repos []ghentity.Repository
	err   error
}

func (m *mockListUserReposUseCase) Execute(_ context.Context, _ usecase.ListUserReposInput) ([]ghentity.Repository, error) {
	return m.repos, m.err
}

type mockListUserOrgsUseCase struct {
	orgs []ghentity.Organization
	err  error
}

func (m *mockListUserOrgsUseCase) Execute(_ context.Context, _ usecase.ListUserOrgsInput) ([]ghentity.Organization, error) {
	return m.orgs, m.err
}

type mockListOrgReposUseCase struct {
	repos []ghentity.Repository
	err   error
}

func (m *mockListOrgReposUseCase) Execute(_ context.Context, _ usecase.ListOrgReposInput) ([]ghentity.Repository, error) {
	return m.repos, m.err
}

type handlerWithMocks struct {
	*Handler
	listUserRepos *mockListUserReposUseCase
	listUserOrgs  *mockListUserOrgsUseCase
	listOrgRepos  *mockListOrgReposUseCase
}

func newHandlerWithMocks(
	listUserRepos *mockListUserReposUseCase,
	listUserOrgs *mockListUserOrgsUseCase,
	listOrgRepos *mockListOrgReposUseCase,
) *handlerWithMocks {
	return &handlerWithMocks{
		Handler: &Handler{
			logger: logger.New(),
		},
		listUserRepos: listUserRepos,
		listUserOrgs:  listUserOrgs,
		listOrgRepos:  listOrgRepos,
	}
}

func (h *handlerWithMocks) GetUserGitHubRepositories(ctx context.Context, request api.GetUserGitHubRepositoriesRequestObject) (api.GetUserGitHubRepositoriesResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserGitHubRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	repos, err := h.listUserRepos.Execute(ctx, usecase.ListUserReposInput{
		Refresh: refresh,
		UserID:  userID,
	})
	if err != nil {
		return h.handleReposError(ctx, err, "list user repositories")
	}

	return api.GetUserGitHubRepositories200JSONResponse{
		Data: mapper.ToGitHubRepositories(repos),
	}, nil
}

func (h *handlerWithMocks) GetUserGitHubOrganizations(ctx context.Context, request api.GetUserGitHubOrganizationsRequestObject) (api.GetUserGitHubOrganizationsResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserGitHubOrganizations401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	orgs, err := h.listUserOrgs.Execute(ctx, usecase.ListUserOrgsInput{
		Refresh: refresh,
		UserID:  userID,
	})
	if err != nil {
		return h.handleOrgsError(ctx, err, "list user organizations")
	}

	return api.GetUserGitHubOrganizations200JSONResponse{
		Data: mapper.ToGitHubOrganizations(orgs),
	}, nil
}

func (h *handlerWithMocks) GetOrganizationRepositories(ctx context.Context, request api.GetOrganizationRepositoriesRequestObject) (api.GetOrganizationRepositoriesResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetOrganizationRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	repos, err := h.listOrgRepos.Execute(ctx, usecase.ListOrgReposInput{
		OrgLogin: request.Org,
		Refresh:  refresh,
		UserID:   userID,
	})
	if err != nil {
		return h.handleOrgReposError(ctx, err, "list organization repositories")
	}

	return api.GetOrganizationRepositories200JSONResponse{
		Data: mapper.ToGitHubRepositories(repos),
	}, nil
}

func TestHandler_GetUserGitHubRepositories_Unauthorized(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{},
	)

	resp, err := h.GetUserGitHubRepositories(context.Background(), api.GetUserGitHubRepositoriesRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.GetUserGitHubRepositories401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestHandler_GetUserGitHubRepositories_Success(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{
			repos: []ghentity.Repository{
				{ID: 1, Name: "repo1", FullName: "user/repo1", Owner: "user"},
			},
		},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetUserGitHubRepositories(ctx, api.GetUserGitHubRepositoriesRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	successResp, ok := resp.(api.GetUserGitHubRepositories200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}

	if len(successResp.Data) != 1 {
		t.Errorf("expected 1 repository, got %d", len(successResp.Data))
	}
}

func TestHandler_GetUserGitHubRepositories_NoToken(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{err: domain.ErrNoGitHubToken},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetUserGitHubRepositories(ctx, api.GetUserGitHubRepositoriesRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.GetUserGitHubRepositories401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestHandler_GetUserGitHubRepositories_RateLimited(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{err: &domain.RateLimitError{Limit: 5000, Remaining: 0}},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetUserGitHubRepositories(ctx, api.GetUserGitHubRepositoriesRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.GetUserGitHubRepositories429ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 429 response, got %T", resp)
	}
}

func TestHandler_GetUserGitHubOrganizations_Success(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{},
		&mockListUserOrgsUseCase{
			orgs: []ghentity.Organization{
				{ID: 1, Login: "org1"},
			},
		},
		&mockListOrgReposUseCase{},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetUserGitHubOrganizations(ctx, api.GetUserGitHubOrganizationsRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	successResp, ok := resp.(api.GetUserGitHubOrganizations200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}

	if len(successResp.Data) != 1 {
		t.Errorf("expected 1 organization, got %d", len(successResp.Data))
	}
}

func TestHandler_GetOrganizationRepositories_Success(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{
			repos: []ghentity.Repository{
				{ID: 1, Name: "repo1", FullName: "org/repo1", Owner: "org"},
			},
		},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetOrganizationRepositories(ctx, api.GetOrganizationRepositoriesRequestObject{Org: "org"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	successResp, ok := resp.(api.GetOrganizationRepositories200JSONResponse)
	if !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}

	if len(successResp.Data) != 1 {
		t.Errorf("expected 1 repository, got %d", len(successResp.Data))
	}
}

func TestHandler_GetOrganizationRepositories_NotFound(t *testing.T) {
	h := newHandlerWithMocks(
		&mockListUserReposUseCase{},
		&mockListUserOrgsUseCase{},
		&mockListOrgReposUseCase{err: domain.ErrOrganizationNotFound},
	)

	ctx := contextWithUserID("user-123")
	resp, err := h.GetOrganizationRepositories(ctx, api.GetOrganizationRepositoriesRequestObject{Org: "unknown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.GetOrganizationRepositories404ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 404 response, got %T", resp)
	}
}

func TestNewHandler_Validation(t *testing.T) {
	listUserRepos := usecase.NewListUserReposUseCase(nil, nil, nil)
	listUserOrgs := usecase.NewListUserOrgsUseCase(nil, nil, nil, nil)
	listOrgRepos := usecase.NewListOrgReposUseCase(nil, nil, nil, nil, nil)

	tests := []struct {
		name    string
		cfg     *HandlerConfig
		wantErr bool
	}{
		{"nil config", nil, true},
		{"nil logger", &HandlerConfig{
			ListUserRepos: listUserRepos,
			ListUserOrgs:  listUserOrgs,
			ListOrgRepos:  listOrgRepos,
		}, true},
		{"nil listUserRepos", &HandlerConfig{
			Logger:       logger.New(),
			ListUserOrgs: listUserOrgs,
			ListOrgRepos: listOrgRepos,
		}, true},
		{"nil listUserOrgs", &HandlerConfig{
			Logger:        logger.New(),
			ListUserRepos: listUserRepos,
			ListOrgRepos:  listOrgRepos,
		}, true},
		{"nil listOrgRepos", &HandlerConfig{
			Logger:        logger.New(),
			ListUserRepos: listUserRepos,
			ListUserOrgs:  listUserOrgs,
		}, true},
		{"valid", &HandlerConfig{
			Logger:        logger.New(),
			ListUserRepos: listUserRepos,
			ListUserOrgs:  listUserOrgs,
			ListOrgRepos:  listOrgRepos,
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHandler(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
