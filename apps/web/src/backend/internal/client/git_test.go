package client

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestGetLatestCommitSHA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("handles public repository", func(t *testing.T) {
		client := NewGitClient()
		ctx := context.Background()

		sha, err := client.GetLatestCommitSHA(ctx, "torvalds", "linux")

		if err != nil {
			if errors.Is(err, ErrRepoNotFound) || errors.Is(err, ErrForbidden) || errors.Is(err, ErrInvalidResponse) {
				t.Log("Known error type returned:", err)
			} else {
				t.Log("Error returned:", err)
			}
		} else {
			if sha == "" {
				t.Error("expected non-empty SHA")
			}
			t.Log("Successfully retrieved SHA from public repo:", sha[:8])
		}
	})

	t.Run("returns error for non-existent repository", func(t *testing.T) {
		client := NewGitClient()
		ctx := context.Background()

		_, err := client.GetLatestCommitSHA(ctx, "nonexistent-user-12345", "nonexistent-repo-67890")

		if err == nil {
			t.Error("expected error for non-existent repository")
		}

		if !errors.Is(err, ErrRepoNotFound) && !errors.Is(err, ErrInvalidResponse) {
			t.Logf("expected ErrRepoNotFound or ErrInvalidResponse, got: %v", err)
		}
	})

	t.Run("returns error for private repository without token", func(t *testing.T) {
		client := NewGitClient()
		ctx := context.Background()

		// This should fail because credential is disabled
		_, err := client.GetLatestCommitSHA(ctx, "KubrickCode", "notag")

		if err == nil {
			t.Error("expected error for private repository without token")
		}
	})
}

func TestGetLatestCommitSHAWithToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("handles public repository with token", func(t *testing.T) {
		client := NewGitClient()
		ctx := context.Background()

		// Token doesn't need to be valid for public repos
		sha, err := client.GetLatestCommitSHAWithToken(ctx, "torvalds", "linux", "dummy-token")

		if err != nil {
			t.Log("Error returned:", err)
		} else {
			if sha == "" {
				t.Error("expected non-empty SHA")
			}
			t.Log("Successfully retrieved SHA:", sha[:8])
		}
	})
}
