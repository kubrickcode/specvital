package middleware

import "net/http"

const (
	headerContentTypeOptions = "X-Content-Type-Options"
	headerFrameOptions       = "X-Frame-Options"
	headerXSSProtection      = "X-XSS-Protection"
	headerReferrerPolicy     = "Referrer-Policy"
	headerCSPReportOnly      = "Content-Security-Policy-Report-Only"

	valueNoSniff        = "nosniff"
	valueDeny           = "DENY"
	valueXSSBlock       = "1; mode=block"
	valueStrictOrigin   = "strict-origin-when-cross-origin"
	valueDefaultSrcSelf = "default-src 'self'"
)

// SecurityHeaders sets security headers. Currently using CSP Report-Only mode.
// TODO: Switch to enforcing CSP after monitoring violations and adding connect-src for frontend origins.
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set(headerContentTypeOptions, valueNoSniff)
			h.Set(headerFrameOptions, valueDeny)
			h.Set(headerXSSProtection, valueXSSBlock)
			h.Set(headerReferrerPolicy, valueStrictOrigin)
			h.Set(headerCSPReportOnly, valueDefaultSrcSelf)

			next.ServeHTTP(w, r)
		})
	}
}
