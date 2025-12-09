package source

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewGitSource(t *testing.T) {
	if !isGitInstalled() {
		t.Skip("git not installed")
	}

	t.Run("should clone local git repository", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		// When
		src, err := NewGitSource(ctx, repoDir, nil)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer src.Close()

		if src.Root() == "" {
			t.Error("expected non-empty root path")
		}
		if src.Root() == repoDir {
			t.Error("expected cloned path to differ from source repo")
		}
	})

	t.Run("should read file from cloned repository", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}
		defer src.Close()

		// When
		reader, err := src.Open(ctx, "test.txt")

		// Then
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read content: %v", err)
		}
		if string(content) != "hello git" {
			t.Errorf("expected 'hello git', got %q", content)
		}
	})

	t.Run("should clone specific branch", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepoWithBranch(t, "feature")
		ctx := context.Background()
		opts := &GitOptions{Branch: "feature"}

		// When
		src, err := NewGitSource(ctx, repoDir, opts)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer src.Close()

		reader, err := src.Open(ctx, "feature.txt")
		if err != nil {
			t.Fatalf("failed to open feature file: %v", err)
		}
		defer reader.Close()

		content, _ := io.ReadAll(reader)
		if string(content) != "feature branch" {
			t.Errorf("expected 'feature branch', got %q", content)
		}
	})

	t.Run("should fail with invalid URL", func(t *testing.T) {
		// Given
		ctx := context.Background()
		invalidURL := "://invalid-url"

		// When
		src, err := NewGitSource(ctx, invalidURL, nil)

		// Then
		if err == nil {
			src.Close()
			t.Fatal("expected error for invalid URL")
		}
		if !errors.Is(err, ErrInvalidPath) && !errors.Is(err, ErrGitCloneFailed) {
			t.Errorf("expected ErrInvalidPath or ErrGitCloneFailed, got: %v", err)
		}
	})

	t.Run("should fail with non-existent repository", func(t *testing.T) {
		// Given
		ctx := context.Background()
		nonExistentPath := "/nonexistent/repo/path"

		// When
		src, err := NewGitSource(ctx, nonExistentPath, nil)

		// Then
		if err == nil {
			src.Close()
			t.Fatal("expected error for non-existent repository")
		}
		if !errors.Is(err, ErrGitCloneFailed) {
			t.Errorf("expected ErrGitCloneFailed, got: %v", err)
		}
	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// When
		src, err := NewGitSource(ctx, repoDir, nil)

		// Then
		if err == nil {
			src.Close()
			t.Fatal("expected error for cancelled context")
		}
	})

	t.Run("should respect context timeout", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(10 * time.Millisecond)

		// When
		src, err := NewGitSource(ctx, repoDir, nil)

		// Then
		if err == nil {
			src.Close()
			t.Fatal("expected error for timed out context")
		}
	})
}

func TestGitSource_CommitSHA(t *testing.T) {
	if !isGitInstalled() {
		t.Skip("git not installed")
	}

	t.Run("should return valid commit SHA", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}
		defer src.Close()

		// When
		sha := src.CommitSHA()

		// Then
		if sha == "" {
			t.Error("expected non-empty commit SHA")
		}
		if len(sha) != 40 {
			t.Errorf("expected 40 character SHA, got %d: %q", len(sha), sha)
		}
	})

	t.Run("should return consistent SHA on multiple calls", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}
		defer src.Close()

		// When
		sha1 := src.CommitSHA()
		sha2 := src.CommitSHA()
		sha3 := src.CommitSHA()

		// Then
		if sha1 != sha2 || sha2 != sha3 {
			t.Errorf("expected consistent SHA, got %q, %q, %q", sha1, sha2, sha3)
		}
	})
}

func TestGitSource_Close(t *testing.T) {
	if !isGitInstalled() {
		t.Skip("git not installed")
	}

	t.Run("should remove cloned directory on close", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}

		clonedPath := src.Root()
		if _, err := os.Stat(clonedPath); os.IsNotExist(err) {
			t.Fatal("cloned directory should exist before close")
		}

		// When
		err = src.Close()

		// Then
		if err != nil {
			t.Errorf("close returned error: %v", err)
		}
		if _, err := os.Stat(clonedPath); !os.IsNotExist(err) {
			t.Error("cloned directory should be removed after close")
		}
	})

	t.Run("should be idempotent", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}

		// When
		err1 := src.Close()
		err2 := src.Close()
		err3 := src.Close()

		// Then
		if err1 != nil {
			t.Errorf("first close returned error: %v", err1)
		}
		if err2 != nil {
			t.Errorf("second close returned error: %v", err2)
		}
		if err3 != nil {
			t.Errorf("third close returned error: %v", err3)
		}
	})

	t.Run("should be safe for concurrent calls", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		src, err := NewGitSource(ctx, repoDir, nil)
		if err != nil {
			t.Fatalf("failed to create git source: %v", err)
		}

		// When
		const goroutines = 10
		errChan := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				errChan <- src.Close()
			}()
		}

		// Then
		for i := 0; i < goroutines; i++ {
			if err := <-errChan; err != nil {
				t.Errorf("concurrent close returned error: %v", err)
			}
		}
	})
}

func TestVerifyGitInstalled(t *testing.T) {
	t.Run("should succeed when git is installed", func(t *testing.T) {
		if !isGitInstalled() {
			t.Skip("git not installed")
		}

		// When
		err := VerifyGitInstalled()

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestVerifyRepository(t *testing.T) {
	if !isGitInstalled() {
		t.Skip("git not installed")
	}

	t.Run("should verify local repository", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()

		// When
		err := VerifyRepository(ctx, repoDir, nil)

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should fail for non-existent repository", func(t *testing.T) {
		// Given
		ctx := context.Background()
		nonExistentPath := "/nonexistent/repo/path"

		// When
		err := VerifyRepository(ctx, nonExistentPath, nil)

		// Then
		if err == nil {
			t.Fatal("expected error for non-existent repository")
		}
		if !errors.Is(err, ErrRepositoryNotFound) {
			t.Errorf("expected ErrRepositoryNotFound, got: %v", err)
		}
	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// When
		err := VerifyRepository(ctx, repoDir, nil)

		// Then
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})
}

func TestInjectCredentials(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		creds    *GitCredentials
		wantURL  string
		wantErr  bool
	}{
		{
			name:    "should not modify URL without credentials",
			repoURL: "https://github.com/owner/repo.git",
			creds:   nil,
			wantURL: "https://github.com/owner/repo.git",
		},
		{
			name:    "should inject username and password",
			repoURL: "https://github.com/owner/repo.git",
			creds:   &GitCredentials{Username: "user", Password: "pass"},
			wantURL: "https://user:pass@github.com/owner/repo.git",
		},
		{
			name:    "should use oauth2 as default username for token",
			repoURL: "https://github.com/owner/repo.git",
			creds:   &GitCredentials{Password: "token123"},
			wantURL: "https://oauth2:token123@github.com/owner/repo.git",
		},
		{
			name:    "should inject username only",
			repoURL: "https://github.com/owner/repo.git",
			creds:   &GitCredentials{Username: "user"},
			wantURL: "https://user@github.com/owner/repo.git",
		},
		{
			name:    "should handle empty credentials",
			repoURL: "https://github.com/owner/repo.git",
			creds:   &GitCredentials{},
			wantURL: "https://github.com/owner/repo.git",
		},
		{
			name:    "should fail for invalid URL",
			repoURL: "://invalid",
			creds:   &GitCredentials{Username: "user", Password: "pass"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			gotURL, err := injectCredentials(tt.repoURL, tt.creds)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("injectCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotURL != tt.wantURL {
				t.Errorf("injectCredentials() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestSanitizeOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
		creds  *GitCredentials
		want   string
	}{
		{
			name:   "should not modify output without credentials",
			output: "error: repository not found",
			creds:  nil,
			want:   "error: repository not found",
		},
		{
			name:   "should redact username",
			output: "fatal: Authentication failed for 'https://myuser@github.com/repo'",
			creds:  &GitCredentials{Username: "myuser"},
			want:   "fatal: Authentication failed for 'https://[REDACTED]@github.com/repo'",
		},
		{
			name:   "should redact password",
			output: "fatal: could not read Password for 'https://github.com': secret123",
			creds:  &GitCredentials{Password: "secret123"},
			want:   "fatal: could not read Password for 'https://github.com': [REDACTED]",
		},
		{
			name:   "should redact URL-encoded password",
			output: "https://user:my%40pass@github.com",
			creds:  &GitCredentials{Password: "my@pass"},
			want:   "https://user:[REDACTED]@github.com",
		},
		{
			name:   "should redact both username and password",
			output: "error: myuser:secret123@github.com",
			creds:  &GitCredentials{Username: "myuser", Password: "secret123"},
			want:   "error: [REDACTED]:[REDACTED]@github.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got := sanitizeOutput(tt.output, tt.creds)

			// Then
			if got != tt.want {
				t.Errorf("sanitizeOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoPathFromURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		want    string
		wantErr bool
	}{
		{
			name:    "should extract path from HTTPS URL",
			repoURL: "https://github.com/owner/repo.git",
			want:    filepath.Join("github.com", "owner", "repo"),
		},
		{
			name:    "should extract path without .git suffix",
			repoURL: "https://github.com/owner/repo",
			want:    filepath.Join("github.com", "owner", "repo"),
		},
		{
			name:    "should handle GitLab URL",
			repoURL: "https://gitlab.com/group/subgroup/repo.git",
			want:    filepath.Join("gitlab.com", "group", "subgroup", "repo"),
		},
		{
			name:    "should strip credentials from URL",
			repoURL: "https://user:secret@github.com/owner/repo.git",
			want:    filepath.Join("github.com", "owner", "repo"),
		},
		{
			name:    "should fail for invalid URL",
			repoURL: "://invalid",
			wantErr: true,
		},
		{
			name:    "should fail for URL without host",
			repoURL: "/owner/repo",
			wantErr: true,
		},
		{
			name:    "should fail for URL without path",
			repoURL: "https://github.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := RepoPathFromURL(tt.repoURL)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoPathFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RepoPathFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitSource_ImplementsSource(t *testing.T) {
	t.Run("should implement Source interface", func(t *testing.T) {
		var _ Source = (*GitSource)(nil)
	})
}

func TestRepoPathFromURL_CredentialSafety(t *testing.T) {
	t.Run("should not leak credentials in error messages", func(t *testing.T) {
		// Given
		urlWithCreds := "https://secretuser:secretpass@github.com"

		// When
		_, err := RepoPathFromURL(urlWithCreds)

		// Then
		if err == nil {
			t.Fatal("expected error for URL without path")
		}

		errMsg := err.Error()
		if strings.Contains(errMsg, "secretuser") {
			t.Error("username leaked in error message")
		}
		if strings.Contains(errMsg, "secretpass") {
			t.Error("password leaked in error message")
		}
	})
}

// createLocalGitRepo creates a temporary git repository for testing.
func createLocalGitRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v\n%s", args, err, out)
		}
	}

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello git"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	addCommit := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "initial"},
	}

	for _, args := range addCommit {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v\n%s", args, err, out)
		}
	}

	return tmpDir
}

// createLocalGitRepoWithBranch creates a git repository with an additional branch.
func createLocalGitRepoWithBranch(t *testing.T, branchName string) string {
	t.Helper()

	repoDir := createLocalGitRepo(t)

	branchCommands := [][]string{
		{"git", "checkout", "-b", branchName},
	}

	for _, args := range branchCommands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to run %v: %v\n%s", args, err, out)
		}
	}

	featureFile := filepath.Join(repoDir, "feature.txt")
	if err := os.WriteFile(featureFile, []byte("feature branch"), 0644); err != nil {
		t.Fatalf("failed to create feature file: %v", err)
	}

	commitCommands := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "feature commit"},
		{"git", "checkout", "master"},
	}

	for _, args := range commitCommands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			if strings.Contains(string(out), "error: pathspec 'master' did not match") {
				cmd = exec.Command("git", "checkout", "main")
				cmd.Dir = repoDir
				cmd.CombinedOutput()
			} else {
				t.Fatalf("failed to run %v: %v\n%s", args, err, out)
			}
		}
	}

	return repoDir
}

func isGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{
			name:    "should accept valid branch name",
			branch:  "feature/new-feature",
			wantErr: false,
		},
		{
			name:    "should accept simple branch name",
			branch:  "main",
			wantErr: false,
		},
		{
			name:    "should reject branch starting with dash",
			branch:  "-malicious",
			wantErr: true,
		},
		{
			name:    "should reject command injection attempt",
			branch:  "--upload-pack=evil",
			wantErr: true,
		},
		{
			name:    "should reject null byte in branch name",
			branch:  "branch\x00name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			err := validateBranchName(tt.branch)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBranchName(%q) error = %v, wantErr %v", tt.branch, err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !errors.Is(err, ErrInvalidPath) {
				t.Errorf("validateBranchName(%q) should return ErrInvalidPath, got %v", tt.branch, err)
			}
		})
	}
}

func TestNewGitSource_MaliciousBranch(t *testing.T) {
	if !isGitInstalled() {
		t.Skip("git not installed")
	}

	t.Run("should reject malicious branch name", func(t *testing.T) {
		// Given
		repoDir := createLocalGitRepo(t)
		ctx := context.Background()
		opts := &GitOptions{Branch: "--upload-pack=malicious"}

		// When
		src, err := NewGitSource(ctx, repoDir, opts)

		// Then
		if err == nil {
			src.Close()
			t.Fatal("expected error for malicious branch name")
		}
		if !errors.Is(err, ErrInvalidPath) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})
}

func TestSanitizeOutput_ShortCredentials(t *testing.T) {
	t.Run("should not redact very short credentials", func(t *testing.T) {
		// Given - short credentials that might match common words
		output := "error: ab cd"
		creds := &GitCredentials{Username: "ab", Password: "cd"}

		// When
		got := sanitizeOutput(output, creds)

		// Then - short credentials should not be redacted to avoid false positives
		if got != output {
			t.Errorf("sanitizeOutput() should not redact short credentials, got %q", got)
		}
	})
}
