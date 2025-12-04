package github

type FileInfo struct {
	Path string `json:"path"`
	SHA  string `json:"sha"`
	Size int    `json:"size,omitempty"`
	Type string `json:"type"` // "file" or "dir"
}

type CommitInfo struct {
	Message string `json:"message"`
	SHA     string `json:"sha"`
}

type RateLimitInfo struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	ResetAt   int64 `json:"resetAt"`
}

type treeResponse struct {
	SHA       string      `json:"sha"`
	Tree      []treeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

type treeEntry struct {
	Mode string `json:"mode"`
	Path string `json:"path"`
	SHA  string `json:"sha"`
	Size int    `json:"size,omitempty"`
	Type string `json:"type"` // "blob" or "tree"
}

type contentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	SHA      string `json:"sha"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
}

type branchResponse struct {
	Commit commitRef `json:"commit"`
	Name   string    `json:"name"`
}

type commitRef struct {
	Commit commitDetail `json:"commit"`
	SHA    string       `json:"sha"`
}

type commitDetail struct {
	Message string `json:"message"`
}

type errorResponse struct {
	DocumentationURL string `json:"documentation_url"`
	Message          string `json:"message"`
}
