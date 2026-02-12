package port

import (
	"context"
	"time"
)

type Repository interface {
	DeleteOrgRepositories(ctx context.Context, userID, orgID string) error
	DeleteUserOrganizations(ctx context.Context, userID string) error
	DeleteUserRepositories(ctx context.Context, userID string) error
	GetOrgIDByLogin(ctx context.Context, login string) (string, error)
	GetOrgRepositories(ctx context.Context, userID, orgID string) ([]RepositoryRecord, error)
	GetUserOrganizations(ctx context.Context, userID string) ([]OrganizationRecord, error)
	GetUserRepositories(ctx context.Context, userID string) ([]RepositoryRecord, error)
	HasOrgRepositories(ctx context.Context, userID, orgID string) (bool, error)
	HasUserOrganizations(ctx context.Context, userID string) (bool, error)
	HasUserRepositories(ctx context.Context, userID string) (bool, error)
	UpsertOrgRepositories(ctx context.Context, userID, orgID string, repos []RepositoryRecord) error
	UpsertUserOrganizations(ctx context.Context, userID string, orgs []OrganizationRecord) error
	UpsertUserRepositories(ctx context.Context, userID string, repos []RepositoryRecord) error
}

type RepositoryRecord struct {
	Archived      bool
	DefaultBranch string
	Description   string
	Disabled      bool
	Fork          bool
	FullName      string
	HTMLURL       string
	ID            int64
	Language      string
	Name          string
	Owner         string
	Private       bool
	PushedAt      *time.Time
	StarCount     int
	Visibility    string
}

type OrganizationRecord struct {
	AvatarURL   string
	Description string
	HTMLURL     string
	ID          int64
	Login       string
	OrgID       string
	Role        string
}
