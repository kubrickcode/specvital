package entity

type AccessStatus string

const (
	AccessStatusAccessible AccessStatus = "accessible"
	AccessStatusPending    AccessStatus = "pending"
	AccessStatusRestricted AccessStatus = "restricted"
)

type Organization struct {
	AccessStatus AccessStatus
	AvatarURL    string
	Description  string
	HTMLURL      string
	ID           int64
	Login        string
	Role         string
}
