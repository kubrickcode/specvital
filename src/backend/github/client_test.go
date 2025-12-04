package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListFiles(t *testing.T) {
	t.Run("should return files from tree API", func(t *testing.T) {
		// Given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "/git/trees/HEAD") {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.URL.Query().Get("recursive") != "1" {
				t.Error("expected recursive=1 query param")
			}

			w.Header().Set("X-RateLimit-Remaining", "4999")
			w.Header().Set("X-RateLimit-Limit", "5000")
			resp := treeResponse{
				SHA: "abc123",
				Tree: []treeEntry{
					{Path: "src/main.go", Type: "blob", SHA: "sha1", Size: 100},
					{Path: "src", Type: "tree", SHA: "sha2"},
					{Path: "README.md", Type: "blob", SHA: "sha3", Size: 50},
				},
				Truncated: false,
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		// When
		files, err := client.ListFiles(context.Background(), "owner", "repo")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(files) != 3 {
			t.Fatalf("expected 3 files, got %d", len(files))
		}
		if files[0].Path != "src/main.go" || files[0].Type != "file" {
			t.Errorf("unexpected first file: %+v", files[0])
		}
		if files[1].Type != "dir" {
			t.Errorf("expected dir type for src folder, got %s", files[1].Type)
		}
	})

	t.Run("should return ErrNotFound for 404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.ListFiles(context.Background(), "owner", "notfound")

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("should return ErrForbidden for private repo", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RateLimit-Remaining", "100")
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.ListFiles(context.Background(), "owner", "private")

		if err != ErrForbidden {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("should return ErrTreeTruncated for large repos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := treeResponse{
				SHA:       "abc123",
				Tree:      []treeEntry{},
				Truncated: true,
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.ListFiles(context.Background(), "owner", "huge-repo")

		if err != ErrTreeTruncated {
			t.Errorf("expected ErrTreeTruncated, got %v", err)
		}
	})

	t.Run("should return error for empty owner", func(t *testing.T) {
		client := NewClient("")

		_, err := client.ListFiles(context.Background(), "", "repo")

		if err != ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("should return error for empty repo", func(t *testing.T) {
		client := NewClient("")

		_, err := client.ListFiles(context.Background(), "owner", "")

		if err != ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestGetFileContent(t *testing.T) {
	t.Run("should decode base64 content", func(t *testing.T) {
		// Given
		expectedContent := "package main\n\nfunc main() {}\n"
		encoded := base64.StdEncoding.EncodeToString([]byte(expectedContent))

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := contentResponse{
				Content:  encoded,
				Encoding: "base64",
				Name:     "main.go",
				Path:     "src/main.go",
				SHA:      "sha123",
				Size:     len(expectedContent),
				Type:     "file",
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		// When
		content, err := client.GetFileContent(context.Background(), "owner", "repo", "src/main.go")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if content != expectedContent {
			t.Errorf("expected %q, got %q", expectedContent, content)
		}
	})

	t.Run("should return ErrNotFound for missing file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.GetFileContent(context.Background(), "owner", "repo", "nonexistent.go")

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("should return error for unexpected encoding", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := contentResponse{
				Content:  "plain text",
				Encoding: "utf-8",
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.GetFileContent(context.Background(), "owner", "repo", "file.txt")

		if err == nil || !strings.Contains(err.Error(), "unexpected encoding") {
			t.Errorf("expected encoding error, got %v", err)
		}
	})

	t.Run("should return error for empty path", func(t *testing.T) {
		client := NewClient("")

		_, err := client.GetFileContent(context.Background(), "owner", "repo", "")

		if err != ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestGetLatestCommit(t *testing.T) {
	t.Run("should return latest commit info", func(t *testing.T) {
		// Given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/branches/main") {
				resp := branchResponse{
					Name: "main",
					Commit: commitRef{
						SHA: "abc123def456",
						Commit: commitDetail{
							Message: "feat: add new feature",
						},
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// Repo info endpoint
			resp := struct {
				DefaultBranch string `json:"default_branch"`
			}{DefaultBranch: "main"}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		// When
		commit, err := client.GetLatestCommit(context.Background(), "owner", "repo")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if commit.SHA != "abc123def456" {
			t.Errorf("expected SHA abc123def456, got %s", commit.SHA)
		}
		if commit.Message != "feat: add new feature" {
			t.Errorf("unexpected message: %s", commit.Message)
		}
	})

	t.Run("should return ErrNotFound for nonexistent repo", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.GetLatestCommit(context.Background(), "owner", "nonexistent")

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("should return error for empty owner", func(t *testing.T) {
		client := NewClient("")

		_, err := client.GetLatestCommit(context.Background(), "", "repo")

		if err != ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestRateLimitTracking(t *testing.T) {
	t.Run("should track rate limit headers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RateLimit-Remaining", "4500")
			w.Header().Set("X-RateLimit-Limit", "5000")
			w.Header().Set("X-RateLimit-Reset", "1699999999")
			resp := treeResponse{SHA: "abc", Tree: []treeEntry{}}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, _ = client.ListFiles(context.Background(), "owner", "repo")

		rateLimit := client.GetRateLimit()
		if rateLimit.Remaining != 4500 {
			t.Errorf("expected remaining 4500, got %d", rateLimit.Remaining)
		}
		if rateLimit.Limit != 5000 {
			t.Errorf("expected limit 5000, got %d", rateLimit.Limit)
		}
		if rateLimit.ResetAt != 1699999999 {
			t.Errorf("expected reset 1699999999, got %d", rateLimit.ResetAt)
		}
	})

	t.Run("should return ErrRateLimited when limit exceeded", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, err := client.ListFiles(context.Background(), "owner", "repo")

		if err != ErrRateLimited {
			t.Errorf("expected ErrRateLimited, got %v", err)
		}
	})
}

func TestAuthorizationHeader(t *testing.T) {
	t.Run("should include token in Authorization header", func(t *testing.T) {
		var receivedAuth string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			resp := treeResponse{SHA: "abc", Tree: []treeEntry{}}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("test-token-123").WithBaseURL(server.URL)

		_, _ = client.ListFiles(context.Background(), "owner", "repo")

		if receivedAuth != "Bearer test-token-123" {
			t.Errorf("expected 'Bearer test-token-123', got %q", receivedAuth)
		}
	})

	t.Run("should not include Authorization header without token", func(t *testing.T) {
		var receivedAuth string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			resp := treeResponse{SHA: "abc", Tree: []treeEntry{}}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient("").WithBaseURL(server.URL)

		_, _ = client.ListFiles(context.Background(), "owner", "repo")

		if receivedAuth != "" {
			t.Errorf("expected empty Authorization, got %q", receivedAuth)
		}
	})
}
