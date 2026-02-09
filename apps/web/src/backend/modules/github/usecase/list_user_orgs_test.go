package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/web/src/backend/modules/github/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type mockInstallationLookup struct {
	installations    map[int64]*port.InstallationInfo
	orgInstallations []port.OrganizationInstallationInfo
	err              error
}

func (m *mockInstallationLookup) GetInstallationByAccountID(_ context.Context, accountID int64) (*port.InstallationInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.installations == nil {
		return nil, nil
	}
	return m.installations[accountID], nil
}

func (m *mockInstallationLookup) GetInstallationsByAccountIDs(_ context.Context, accountIDs []int64) (map[int64]*port.InstallationInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.installations == nil {
		return make(map[int64]*port.InstallationInfo), nil
	}
	result := make(map[int64]*port.InstallationInfo)
	for _, id := range accountIDs {
		if info, exists := m.installations[id]; exists {
			result[id] = info
		}
	}
	return result, nil
}

func (m *mockInstallationLookup) ListOrganizationsByUserID(_ context.Context, _ string) ([]port.OrganizationInstallationInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.orgInstallations, nil
}

func TestListUserOrgsUseCase_AccessStatus_NoInstallation(t *testing.T) {
	repo := &mockRepository{
		hasOrgs: true,
		orgs:    []port.OrganizationRecord{{ID: 123, Login: "org1"}},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	lookup := &mockInstallationLookup{}

	uc := NewListUserOrgsUseCase(mockClientFactory(ghClient), repo, provider, lookup)

	orgs, err := uc.Execute(context.Background(), ListUserOrgsInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].AccessStatus != entity.AccessStatusRestricted {
		t.Errorf("expected AccessStatusRestricted, got %v", orgs[0].AccessStatus)
	}
}

func TestListUserOrgsUseCase_AccessStatus_WithInstallation(t *testing.T) {
	repo := &mockRepository{
		hasOrgs: true,
		orgs:    []port.OrganizationRecord{{ID: 123, Login: "org1"}},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	lookup := &mockInstallationLookup{
		installations: map[int64]*port.InstallationInfo{
			123: {AccountID: 123, InstallationID: 456, IsSuspended: false},
		},
	}

	uc := NewListUserOrgsUseCase(mockClientFactory(ghClient), repo, provider, lookup)

	orgs, err := uc.Execute(context.Background(), ListUserOrgsInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].AccessStatus != entity.AccessStatusAccessible {
		t.Errorf("expected AccessStatusAccessible, got %v", orgs[0].AccessStatus)
	}
}

func TestListUserOrgsUseCase_AccessStatus_SuspendedInstallation(t *testing.T) {
	repo := &mockRepository{
		hasOrgs: true,
		orgs:    []port.OrganizationRecord{{ID: 123, Login: "org1"}},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	lookup := &mockInstallationLookup{
		installations: map[int64]*port.InstallationInfo{
			123: {AccountID: 123, InstallationID: 456, IsSuspended: true},
		},
	}

	uc := NewListUserOrgsUseCase(mockClientFactory(ghClient), repo, provider, lookup)

	orgs, err := uc.Execute(context.Background(), ListUserOrgsInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 1 {
		t.Fatalf("expected 1 org, got %d", len(orgs))
	}
	if orgs[0].AccessStatus != entity.AccessStatusPending {
		t.Errorf("expected AccessStatusPending, got %v", orgs[0].AccessStatus)
	}
}

func TestListUserOrgsUseCase_AccessStatus_MixedStatus(t *testing.T) {
	repo := &mockRepository{
		hasOrgs: true,
		orgs: []port.OrganizationRecord{
			{ID: 100, Login: "org-accessible"},
			{ID: 200, Login: "org-suspended"},
			{ID: 300, Login: "org-restricted"},
		},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	lookup := &mockInstallationLookup{
		installations: map[int64]*port.InstallationInfo{
			100: {AccountID: 100, InstallationID: 1001, IsSuspended: false},
			200: {AccountID: 200, InstallationID: 2002, IsSuspended: true},
		},
	}

	uc := NewListUserOrgsUseCase(mockClientFactory(ghClient), repo, provider, lookup)

	orgs, err := uc.Execute(context.Background(), ListUserOrgsInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 3 {
		t.Fatalf("expected 3 orgs, got %d", len(orgs))
	}

	statusMap := make(map[string]entity.AccessStatus)
	for _, org := range orgs {
		statusMap[org.Login] = org.AccessStatus
	}

	if statusMap["org-accessible"] != entity.AccessStatusAccessible {
		t.Errorf("org-accessible: expected AccessStatusAccessible, got %v", statusMap["org-accessible"])
	}
	if statusMap["org-suspended"] != entity.AccessStatusPending {
		t.Errorf("org-suspended: expected AccessStatusPending, got %v", statusMap["org-suspended"])
	}
	if statusMap["org-restricted"] != entity.AccessStatusRestricted {
		t.Errorf("org-restricted: expected AccessStatusRestricted, got %v", statusMap["org-restricted"])
	}
}

func TestListUserOrgsUseCase_AccessStatus_LookupError(t *testing.T) {
	repo := &mockRepository{
		hasOrgs: true,
		orgs: []port.OrganizationRecord{
			{ID: 100, Login: "org1"},
			{ID: 200, Login: "org2"},
		},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	lookup := &mockInstallationLookup{
		err: errors.New("database connection failed"),
	}

	uc := NewListUserOrgsUseCase(mockClientFactory(ghClient), repo, provider, lookup)

	orgs, err := uc.Execute(context.Background(), ListUserOrgsInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgs) != 2 {
		t.Fatalf("expected 2 orgs, got %d", len(orgs))
	}

	for _, org := range orgs {
		if org.AccessStatus != entity.AccessStatusRestricted {
			t.Errorf("org %s: expected AccessStatusRestricted on lookup error, got %v", org.Login, org.AccessStatus)
		}
	}
}
