package entity

import "time"

type AccountType string

const (
	AccountTypeOrganization AccountType = "organization"
	AccountTypeUser         AccountType = "user"
)

type Installation struct {
	AccountAvatarURL *string
	AccountID        int64
	AccountLogin     string
	AccountType      AccountType
	CreatedAt        time.Time
	ID               string
	InstallationID   int64
	InstallerUserID  *string
	SuspendedAt      *time.Time
	UpdatedAt        time.Time
}

func (i *Installation) IsOrganization() bool {
	return i.AccountType == AccountTypeOrganization
}

func (i *Installation) IsSuspended() bool {
	return i.SuspendedAt != nil
}
