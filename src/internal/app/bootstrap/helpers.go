package bootstrap

import (
	"fmt"
	"net/url"
)

const defaultConcurrency = 5

// maskURL returns a sanitized URL for logging (hides credentials).
func maskURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "[invalid-url]"
	}

	host := parsed.Host
	if len(host) > 30 {
		host = host[:30] + "..."
	}

	userPart := ""
	if parsed.User != nil {
		userPart = parsed.User.Username() + ":****@"
	}

	return fmt.Sprintf("%s://%s%s/...", parsed.Scheme, userPart, host)
}
