package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestGetLatestCommitSHAWithToken(t *testing.T) {
	t.Run("returns SHA when API succeeds", func(t *testing.T) {
		// Given: mock GitHub API server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify headers
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Error("expected Authorization header with Bearer token")
			}
			if r.Header.Get("Accept") != "application/vnd.github+json" {
				t.Error("expected Accept header")
			}
			if r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
				t.Error("expected X-GitHub-Api-Version header")
			}

			// Verify path
			expectedPath := "/repos/owner/repo/commits/HEAD"
			if r.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
			}

			response := map[string]string{"sha": "abc123def456"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with token
		sha, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "test-token")

		// Then: returns expected SHA
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sha != "abc123def456" {
			t.Errorf("expected sha abc123def456, got %s", sha)
		}
	})

	t.Run("returns ErrRepoNotFound on 404", func(t *testing.T) {
		// Given: mock server returning 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with token
		_, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "test-token")

		// Then: returns ErrRepoNotFound
		if !errors.Is(err, ErrRepoNotFound) {
			t.Errorf("expected ErrRepoNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden on 403", func(t *testing.T) {
		// Given: mock server returning 403
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with token
		_, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "test-token")

		// Then: returns ErrForbidden
		if !errors.Is(err, ErrForbidden) {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns ErrForbidden on 401", func(t *testing.T) {
		// Given: mock server returning 401
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with invalid token
		_, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "invalid-token")

		// Then: returns ErrForbidden
		if !errors.Is(err, ErrForbidden) {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns error on empty SHA", func(t *testing.T) {
		// Given: mock server returning empty SHA
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]string{"sha": ""}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with token
		_, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "test-token")

		// Then: returns ErrInvalidResponse
		if !errors.Is(err, ErrInvalidResponse) {
			t.Errorf("expected ErrInvalidResponse, got %v", err)
		}
	})

	t.Run("returns error on unexpected status code", func(t *testing.T) {
		// Given: mock server returning 500
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal error"))
		}))
		defer server.Close()

		client := NewGitClientWithOptions(http.DefaultClient, server.URL)
		ctx := context.Background()

		// When: calling with token
		_, err := client.GetLatestCommitSHAWithToken(ctx, "owner", "repo", "test-token")

		// Then: returns ErrInvalidResponse
		if !errors.Is(err, ErrInvalidResponse) {
			t.Errorf("expected ErrInvalidResponse, got %v", err)
		}
	})
}

func TestGetLatestCommitSHA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("handles public repository", func(t *testing.T) {
		// Given: git client (integration test - requires git binary)
		client := NewGitClient()
		ctx := context.Background()

		// When: getting commit SHA for known public repo
		_, err := client.GetLatestCommitSHA(ctx, "torvalds", "linux")

		// Then: should succeed or fail gracefully
		if err != nil {
			if errors.Is(err, ErrRepoNotFound) || errors.Is(err, ErrForbidden) || errors.Is(err, ErrInvalidResponse) {
				t.Log("Known error type returned:", err)
			} else {
				t.Log("Error returned:", err)
			}
		} else {
			t.Log("Successfully retrieved SHA from public repo")
		}
	})

	t.Run("returns error for non-existent repository", func(t *testing.T) {
		// Given: git client
		client := NewGitClient()
		ctx := context.Background()

		// When: getting commit SHA for non-existent repo
		_, err := client.GetLatestCommitSHA(ctx, "nonexistent-user-12345", "nonexistent-repo-67890")

		// Then: should return error
		if err == nil {
			t.Error("expected error for non-existent repository")
		}

		if !errors.Is(err, ErrRepoNotFound) && !errors.Is(err, ErrInvalidResponse) {
			t.Logf("expected ErrRepoNotFound or ErrInvalidResponse, got: %v", err)
		}
	})
}

func TestNewGitClientWithOptions(t *testing.T) {
	t.Run("uses defaults when nil/empty provided", func(t *testing.T) {
		client := NewGitClientWithOptions(nil, "")

		// Type assert to check internal state
		gc, ok := client.(*gitClient)
		if !ok {
			t.Fatal("expected *gitClient type")
		}

		if gc.httpClient == nil {
			t.Error("expected non-nil httpClient")
		}
		if gc.baseURL != defaultGitHubAPIURL {
			t.Errorf("expected baseURL %s, got %s", defaultGitHubAPIURL, gc.baseURL)
		}
	})

	t.Run("uses provided values", func(t *testing.T) {
		customClient := &http.Client{}
		customURL := "https://custom.api.com"

		client := NewGitClientWithOptions(customClient, customURL)

		gc, ok := client.(*gitClient)
		if !ok {
			t.Fatal("expected *gitClient type")
		}

		if gc.httpClient != customClient {
			t.Error("expected custom httpClient")
		}
		if gc.baseURL != customURL {
			t.Errorf("expected baseURL %s, got %s", customURL, gc.baseURL)
		}
	})
}
