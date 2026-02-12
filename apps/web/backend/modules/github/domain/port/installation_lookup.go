package port

import "context"

// InstallationInfo contains minimal information about a GitHub App installation.
// This is a cross-module contract - do not add github-app domain specific fields.
type InstallationInfo struct {
	AccountID      int64
	InstallationID int64
	IsSuspended    bool
}

// OrganizationInstallationInfo contains organization info from GitHub App installation.
// Used to supplement OAuth org list when OAuth access is restricted.
type OrganizationInstallationInfo struct {
	AccountAvatarURL *string
	AccountID        int64
	AccountLogin     string
	IsSuspended      bool
}

// InstallationLookup provides access to GitHub App installation information.
// Used by github module to determine access status and token selection.
type InstallationLookup interface {
	// GetInstallationByAccountID returns the installation for a GitHub account.
	// Returns (nil, nil) if no installation exists for the account.
	GetInstallationByAccountID(ctx context.Context, accountID int64) (*InstallationInfo, error)

	// GetInstallationsByAccountIDs returns installations for multiple accounts in a single query.
	// Missing installations are not included in the returned map (no error for missing).
	GetInstallationsByAccountIDs(ctx context.Context, accountIDs []int64) (map[int64]*InstallationInfo, error)

	// ListOrganizationsByUserID returns organizations where user has installed GitHub App.
	// Used to supplement OAuth org list when OAuth access is restricted.
	ListOrganizationsByUserID(ctx context.Context, userID string) ([]OrganizationInstallationInfo, error)
}

// InstallationTokenProvider creates installation access tokens for GitHub API calls.
type InstallationTokenProvider interface {
	// GetInstallationToken creates a new installation token.
	// Tokens are typically valid for 1 hour.
	GetInstallationToken(ctx context.Context, installationID int64) (string, error)
}
