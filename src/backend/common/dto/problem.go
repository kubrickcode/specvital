package dto

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	ContentTypeHeader     = "Content-Type"
	JSONContentType       = "application/json"
	ProblemJSONType       = "application/problem+json"
	ProblemTypeAboutBlank = "about:blank"
)

// ProblemDetail represents RFC 7807 error response format.
type ProblemDetail struct {
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

func SendProblemDetail(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
	problem := ProblemDetail{
		Detail:   detail,
		Instance: r.URL.Path,
		Status:   status,
		Title:    title,
		Type:     ProblemTypeAboutBlank,
	}
	w.Header().Set(ContentTypeHeader, ProblemJSONType)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(problem); err != nil {
		slog.Error("failed to encode problem detail", "error", err)
	}
}
