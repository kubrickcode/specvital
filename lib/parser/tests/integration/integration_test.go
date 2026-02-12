//go:build integration

package integration

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/kubrickcode/specvital/lib/parser"
	"github.com/kubrickcode/specvital/lib/source"

	_ "github.com/kubrickcode/specvital/lib/parser/strategies/cargotest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/cypress"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/gotesting"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/gtest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/jest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/junit4"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/junit5"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/kotest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/minitest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/mocha"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/mstest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/nunit"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/phpunit"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/playwright"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/pytest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/rspec"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/swift-testing"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/testng"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/unittest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/vitest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/xctest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/xunit"
)

const scanTimeout = 10 * time.Minute

var updateSnapshots = flag.Bool("update", false, "update golden snapshots")

func TestSingleFramework(t *testing.T) {
	repos, err := LoadRepos()
	if err != nil {
		t.Fatalf("load repos.yaml: %v", err)
	}

	for _, repo := range repos.Repositories {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			t.Parallel()

			cloneResult, err := CloneRepo(repo)
			if err != nil {
				t.Fatalf("clone %s: %v", repo.Name, err)
			}

			if cloneResult.FromCache {
				t.Logf("using cached repository: %s", cloneResult.Path)
			} else {
				t.Logf("cloned repository: %s", cloneResult.Path)
			}

			ctx, cancel := context.WithTimeout(context.Background(), scanTimeout)
			defer cancel()

			src, err := source.NewLocalSource(cloneResult.Path)
			if err != nil {
				t.Fatalf("create source for %s: %v", repo.Name, err)
			}
			defer src.Close()

			scanResult, err := parser.Scan(ctx, src)
			if err != nil {
				t.Fatalf("scan %s: %v", repo.Name, err)
			}

			t.Logf("scan stats: files=%d, matched=%d, tests=%d, duration=%v",
				scanResult.Stats.FilesScanned,
				scanResult.Stats.FilesMatched,
				scanResult.Inventory.CountTests(),
				scanResult.Stats.Duration,
			)

			// Basic sanity checks
			if scanResult.Stats.FilesMatched == 0 {
				t.Errorf("expected at least 1 matched file, got 0")
			}

			if scanResult.Inventory.CountTests() == 0 {
				t.Errorf("expected at least 1 test, got 0")
			}

			frameworkCount := countByFramework(scanResult)
			validateFrameworkMatch(t, repo.Frameworks, frameworkCount)

			t.Logf("framework distribution: %v", frameworkCount)

			// Snapshot comparison
			actualSnapshot := SnapshotFromScanResult(repo, scanResult, cloneResult.Path)

			if *updateSnapshots {
				if err := SaveSnapshot(actualSnapshot); err != nil {
					t.Fatalf("save snapshot: %v", err)
				}
				t.Logf("updated snapshot for %s", repo.Name)
				return
			}

			// Skip strict snapshot comparison for nondeterministic repositories
			if repo.Nondeterministic {
				t.Logf("skipping strict snapshot comparison for nondeterministic repo %s", repo.Name)
				return
			}

			expectedSnapshot, err := LoadSnapshot(repo.Name, repo.Ref)
			if err != nil {
				t.Fatalf("load snapshot: %v", err)
			}

			diff := CompareSnapshots(expectedSnapshot, actualSnapshot)
			if !diff.IsEmpty() {
				t.Errorf("snapshot mismatch for %s:\n%s", repo.Name, diff.String())
			}
		})
	}
}

func TestScanStreaming(t *testing.T) {
	repos, err := LoadRepos()
	if err != nil {
		t.Fatalf("load repos.yaml: %v", err)
	}

	// Use the first repository for streaming test
	if len(repos.Repositories) == 0 {
		t.Skip("no repositories configured")
	}
	repo := repos.Repositories[0]

	cloneResult, err := CloneRepo(repo)
	if err != nil {
		t.Fatalf("clone %s: %v", repo.Name, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), scanTimeout)
	defer cancel()

	src, err := source.NewLocalSource(cloneResult.Path)
	if err != nil {
		t.Fatalf("create source for %s: %v", repo.Name, err)
	}
	defer src.Close()

	// Test streaming API
	results, err := parser.ScanStreaming(ctx, src)
	if err != nil {
		t.Fatalf("ScanStreaming: %v", err)
	}

	var filesProcessed, filesMatched, parseErrors int
	for result := range results {
		filesProcessed++
		if result.Err != nil {
			parseErrors++
			continue
		}
		if result.File != nil {
			filesMatched++
		}
	}

	t.Logf("streaming stats: processed=%d, matched=%d, errors=%d", filesProcessed, filesMatched, parseErrors)

	if filesMatched == 0 {
		t.Errorf("expected at least 1 matched file via streaming, got 0")
	}

	// Compare with batch Scan() to verify consistency
	scanResult, err := parser.Scan(ctx, src)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if filesMatched != scanResult.Stats.FilesMatched {
		t.Errorf("streaming matched %d files, but Scan matched %d", filesMatched, scanResult.Stats.FilesMatched)
	}
}

func countByFramework(result *parser.ScanResult) map[string]int {
	counts := make(map[string]int)
	for _, file := range result.Inventory.Files {
		counts[file.Framework]++
	}
	return counts
}

// validateFrameworkMatch ensures expected frameworks exactly match actual frameworks.
// This prevents silent failures where secondary frameworks go unvalidated.
func validateFrameworkMatch(t *testing.T, expected []string, actual map[string]int) {
	t.Helper()

	expectedSet := make(map[string]bool)
	for _, fw := range expected {
		expectedSet[fw] = true
	}

	actualSet := make(map[string]bool)
	for fw := range actual {
		actualSet[fw] = true
	}

	var missing, extra []string
	for fw := range expectedSet {
		if !actualSet[fw] {
			missing = append(missing, fw)
		}
	}
	for fw := range actualSet {
		if !expectedSet[fw] {
			extra = append(extra, fw)
		}
	}

	if len(missing) > 0 {
		t.Errorf("expected frameworks not found: %v", missing)
	}
	if len(extra) > 0 {
		t.Errorf("unexpected frameworks detected: %v (update repos.yaml)", extra)
	}
}
