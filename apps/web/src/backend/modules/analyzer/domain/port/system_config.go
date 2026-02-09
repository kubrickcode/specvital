package port

import "context"

type SystemConfigReader interface {
	GetParserVersion(ctx context.Context) (string, error)
}
