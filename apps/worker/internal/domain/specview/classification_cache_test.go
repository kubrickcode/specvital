package specview

import (
	"bytes"
	"testing"
)

func TestGenerateFileSignature(t *testing.T) {
	t.Run("deterministic hash for same files", func(t *testing.T) {
		files := []FileInfo{
			{Path: "test/auth_test.go"},
			{Path: "test/user_test.go"},
		}

		sig1 := GenerateFileSignature(files)
		sig2 := GenerateFileSignature(files)

		if !bytes.Equal(sig1, sig2) {
			t.Error("expected same signature for same files")
		}
	})

	t.Run("different order produces same hash", func(t *testing.T) {
		files1 := []FileInfo{
			{Path: "test/auth_test.go"},
			{Path: "test/user_test.go"},
		}
		files2 := []FileInfo{
			{Path: "test/user_test.go"},
			{Path: "test/auth_test.go"},
		}

		sig1 := GenerateFileSignature(files1)
		sig2 := GenerateFileSignature(files2)

		if !bytes.Equal(sig1, sig2) {
			t.Error("expected same signature regardless of order")
		}
	})

	t.Run("different files produce different hash", func(t *testing.T) {
		files1 := []FileInfo{
			{Path: "test/auth_test.go"},
		}
		files2 := []FileInfo{
			{Path: "test/user_test.go"},
		}

		sig1 := GenerateFileSignature(files1)
		sig2 := GenerateFileSignature(files2)

		if bytes.Equal(sig1, sig2) {
			t.Error("expected different signature for different files")
		}
	})

	t.Run("empty files return valid hash", func(t *testing.T) {
		sig := GenerateFileSignature([]FileInfo{})

		if len(sig) == 0 {
			t.Error("expected non-empty hash for empty files")
		}
		// SHA-256 produces 32 bytes
		if len(sig) != 32 {
			t.Errorf("expected 32-byte hash, got %d bytes", len(sig))
		}
	})

	t.Run("normalizes file paths", func(t *testing.T) {
		files1 := []FileInfo{
			{Path: "test/auth_test.go"},
		}
		files2 := []FileInfo{
			{Path: "test\\auth_test.go"}, // Windows path
		}

		sig1 := GenerateFileSignature(files1)
		sig2 := GenerateFileSignature(files2)

		if !bytes.Equal(sig1, sig2) {
			t.Error("expected same signature after path normalization")
		}
	})
}

func TestTestKey(t *testing.T) {
	t.Run("generates unique key for test identity", func(t *testing.T) {
		key := TestKey("test/auth_test.go", "AuthSuite", "TestLogin")

		if key == "" {
			t.Error("expected non-empty key")
		}
	})

	t.Run("different tests produce different keys", func(t *testing.T) {
		key1 := TestKey("test/auth_test.go", "AuthSuite", "TestLogin")
		key2 := TestKey("test/auth_test.go", "AuthSuite", "TestLogout")

		if key1 == key2 {
			t.Error("expected different keys for different tests")
		}
	})

	t.Run("same test in different suite produces different key", func(t *testing.T) {
		key1 := TestKey("test/auth_test.go", "SuiteA", "TestLogin")
		key2 := TestKey("test/auth_test.go", "SuiteB", "TestLogin")

		if key1 == key2 {
			t.Error("expected different keys for different suites")
		}
	})

	t.Run("same test in different file produces different key", func(t *testing.T) {
		key1 := TestKey("test/auth_test.go", "Suite", "TestLogin")
		key2 := TestKey("test/user_test.go", "Suite", "TestLogin")

		if key1 == key2 {
			t.Error("expected different keys for different files")
		}
	})

	t.Run("normalizes file path in key", func(t *testing.T) {
		key1 := TestKey("test/auth_test.go", "Suite", "TestLogin")
		key2 := TestKey("test\\auth_test.go", "Suite", "TestLogin") // Windows path

		if key1 != key2 {
			t.Error("expected same key after path normalization")
		}
	})

	t.Run("normalizes test name in key", func(t *testing.T) {
		key1 := TestKey("test/auth_test.go", "Suite", "TestLogin")
		key2 := TestKey("test/auth_test.go", "Suite", "  TestLogin  ") // Extra whitespace

		if key1 != key2 {
			t.Error("expected same key after name normalization")
		}
	})
}

func TestBuildTestIndexMap(t *testing.T) {
	t.Run("builds map from Phase1Output", func(t *testing.T) {
		files := []FileInfo{
			{
				Path: "test/auth_test.go",
				Tests: []TestInfo{
					{Index: 0, Name: "TestLogin", SuitePath: "AuthSuite"},
					{Index: 1, Name: "TestLogout", SuitePath: "AuthSuite"},
				},
			},
			{
				Path: "test/user_test.go",
				Tests: []TestInfo{
					{Index: 2, Name: "TestCreateUser", SuitePath: "UserSuite"},
				},
			},
		}

		output := &Phase1Output{
			Domains: []DomainGroup{
				{
					Name: "Auth",
					Features: []FeatureGroup{
						{Name: "Login", TestIndices: []int{0}},
						{Name: "Logout", TestIndices: []int{1}},
					},
				},
				{
					Name: "User",
					Features: []FeatureGroup{
						{Name: "Management", TestIndices: []int{2}},
					},
				},
			},
		}

		indexMap := BuildTestIndexMap(output, files)

		if len(indexMap) != 3 {
			t.Errorf("expected 3 entries, got %d", len(indexMap))
		}

		// Verify test 0 mapping
		key0 := TestKey("test/auth_test.go", "AuthSuite", "TestLogin")
		identity0, ok := indexMap[key0]
		if !ok {
			t.Fatalf("expected key %q to exist", key0)
		}
		if identity0.DomainIndex != 0 {
			t.Errorf("expected DomainIndex 0, got %d", identity0.DomainIndex)
		}
		if identity0.FeatureIndex != 0 {
			t.Errorf("expected FeatureIndex 0, got %d", identity0.FeatureIndex)
		}

		// Verify test 2 mapping
		key2 := TestKey("test/user_test.go", "UserSuite", "TestCreateUser")
		identity2, ok := indexMap[key2]
		if !ok {
			t.Fatalf("expected key %q to exist", key2)
		}
		if identity2.DomainIndex != 1 {
			t.Errorf("expected DomainIndex 1, got %d", identity2.DomainIndex)
		}
		if identity2.FilePath != "test/user_test.go" {
			t.Errorf("expected FilePath 'test/user_test.go', got %q", identity2.FilePath)
		}
	})

	t.Run("returns empty map for nil output", func(t *testing.T) {
		indexMap := BuildTestIndexMap(nil, []FileInfo{})

		if len(indexMap) != 0 {
			t.Errorf("expected empty map, got %d entries", len(indexMap))
		}
	})

	t.Run("returns empty map for empty files", func(t *testing.T) {
		output := &Phase1Output{
			Domains: []DomainGroup{
				{Name: "Domain", Features: []FeatureGroup{{Name: "Feature", TestIndices: []int{0}}}},
			},
		}

		indexMap := BuildTestIndexMap(output, []FileInfo{})

		if len(indexMap) != 0 {
			t.Errorf("expected empty map, got %d entries", len(indexMap))
		}
	})
}
