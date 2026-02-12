package httputil

import (
	"net"
	"net/http"
	"strings"
)

// ExtractClientIP extracts the client IP address from HTTP request headers.
//
// SECURITY NOTE: This function trusts X-Forwarded-For and X-Real-IP headers.
// It assumes the application runs behind a trusted reverse proxy (nginx, cloud LB, etc.)
// that properly sets these headers. If exposed directly to the internet without a proxy,
// clients can spoof their IP address by setting these headers manually.
//
// Header priority:
//  1. X-Forwarded-For (first IP in the list)
//  2. X-Real-IP
//  3. RemoteAddr (direct connection fallback)
func ExtractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
