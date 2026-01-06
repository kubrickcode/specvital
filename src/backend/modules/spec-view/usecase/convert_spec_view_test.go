package usecase

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

func TestGroupUncachedByFile(t *testing.T) {
	codebaseID := uuid.New()

	t.Run("multiple files have independent indices starting from 1", func(t *testing.T) {
		metas := []testMeta{
			{cacheKeyHash: []byte{1}, filePath: "file1.spec.ts", originalName: "test1", suiteHierarchy: "suite1"},
			{cacheKeyHash: []byte{2}, filePath: "file1.spec.ts", originalName: "test2", suiteHierarchy: "suite1"},
			{cacheKeyHash: []byte{3}, filePath: "file2.spec.ts", originalName: "test3", suiteHierarchy: "suite2"},
			{cacheKeyHash: []byte{4}, filePath: "file2.spec.ts", originalName: "test4", suiteHierarchy: "suite2"},
			{cacheKeyHash: []byte{5}, filePath: "file2.spec.ts", originalName: "test5", suiteHierarchy: "suite2"},
		}
		cached := make(map[string]*entity.CacheEntry)

		result := groupUncachedByFile(metas, cached)

		if len(result) != 2 {
			t.Errorf("expected 2 files, got %d", len(result))
		}

		file1 := result["file1.spec.ts"]
		if file1 == nil {
			t.Fatal("file1.spec.ts not found")
		}
		if len(file1.indexToMeta) != 2 {
			t.Errorf("file1 expected 2 tests, got %d", len(file1.indexToMeta))
		}
		if _, ok := file1.indexToMeta["1"]; !ok {
			t.Error("file1 should have index '1'")
		}
		if _, ok := file1.indexToMeta["2"]; !ok {
			t.Error("file1 should have index '2'")
		}

		file2 := result["file2.spec.ts"]
		if file2 == nil {
			t.Fatal("file2.spec.ts not found")
		}
		if len(file2.indexToMeta) != 3 {
			t.Errorf("file2 expected 3 tests, got %d", len(file2.indexToMeta))
		}
		if _, ok := file2.indexToMeta["1"]; !ok {
			t.Error("file2 should have index '1' (not global '3')")
		}
		if _, ok := file2.indexToMeta["2"]; !ok {
			t.Error("file2 should have index '2' (not global '4')")
		}
		if _, ok := file2.indexToMeta["3"]; !ok {
			t.Error("file2 should have index '3' (not global '5')")
		}
	})

	t.Run("index order matches suite order for AI prompt consistency", func(t *testing.T) {
		metas := []testMeta{
			{cacheKeyHash: []byte{1}, filePath: "file.spec.ts", originalName: "testA1", suiteHierarchy: "suiteA"},
			{cacheKeyHash: []byte{2}, filePath: "file.spec.ts", originalName: "testB1", suiteHierarchy: "suiteB"},
			{cacheKeyHash: []byte{3}, filePath: "file.spec.ts", originalName: "testA2", suiteHierarchy: "suiteA"},
			{cacheKeyHash: []byte{4}, filePath: "file.spec.ts", originalName: "testB2", suiteHierarchy: "suiteB"},
		}
		cached := make(map[string]*entity.CacheEntry)

		result := groupUncachedByFile(metas, cached)

		file := result["file.spec.ts"]
		if file == nil {
			t.Fatal("file.spec.ts not found")
		}

		if len(file.suites) != 2 {
			t.Fatalf("expected 2 suites, got %d", len(file.suites))
		}
		if file.suites[0].Hierarchy != "suiteA" {
			t.Errorf("first suite should be suiteA, got %s", file.suites[0].Hierarchy)
		}
		if file.suites[1].Hierarchy != "suiteB" {
			t.Errorf("second suite should be suiteB, got %s", file.suites[1].Hierarchy)
		}

		if file.indexToMeta["1"].originalName != "testA1" {
			t.Errorf("index 1 should be testA1, got %s", file.indexToMeta["1"].originalName)
		}
		if file.indexToMeta["2"].originalName != "testA2" {
			t.Errorf("index 2 should be testA2, got %s", file.indexToMeta["2"].originalName)
		}
		if file.indexToMeta["3"].originalName != "testB1" {
			t.Errorf("index 3 should be testB1, got %s", file.indexToMeta["3"].originalName)
		}
		if file.indexToMeta["4"].originalName != "testB2" {
			t.Errorf("index 4 should be testB2, got %s", file.indexToMeta["4"].originalName)
		}

		if len(file.suites[0].Tests) != 2 || file.suites[0].Tests[0] != "testA1" || file.suites[0].Tests[1] != "testA2" {
			t.Errorf("suiteA tests should be [testA1, testA2], got %v", file.suites[0].Tests)
		}
		if len(file.suites[1].Tests) != 2 || file.suites[1].Tests[0] != "testB1" || file.suites[1].Tests[1] != "testB2" {
			t.Errorf("suiteB tests should be [testB1, testB2], got %v", file.suites[1].Tests)
		}
	})

	t.Run("cached tests are excluded from indexing", func(t *testing.T) {
		metas := []testMeta{
			{cacheKeyHash: []byte{1}, filePath: "file.spec.ts", originalName: "test1", suiteHierarchy: "suite"},
			{cacheKeyHash: []byte{2}, filePath: "file.spec.ts", originalName: "test2", suiteHierarchy: "suite"},
			{cacheKeyHash: []byte{3}, filePath: "file.spec.ts", originalName: "test3", suiteHierarchy: "suite"},
		}
		cached := map[string]*entity.CacheEntry{
			hex.EncodeToString([]byte{2}): {ConvertedName: "cached"},
		}

		result := groupUncachedByFile(metas, cached)

		file := result["file.spec.ts"]
		if file == nil {
			t.Fatal("file.spec.ts not found")
		}
		if len(file.indexToMeta) != 2 {
			t.Errorf("expected 2 uncached tests, got %d", len(file.indexToMeta))
		}
		if file.indexToMeta["1"].originalName != "test1" {
			t.Errorf("index '1' should be test1, got %s", file.indexToMeta["1"].originalName)
		}
		if file.indexToMeta["2"].originalName != "test3" {
			t.Errorf("index '2' should be test3, got %s", file.indexToMeta["2"].originalName)
		}
	})

	t.Run("empty input returns empty result", func(t *testing.T) {
		result := groupUncachedByFile(nil, nil)

		if len(result) != 0 {
			t.Errorf("expected empty result, got %d files", len(result))
		}
	})

	t.Run("all cached returns empty result", func(t *testing.T) {
		metas := []testMeta{
			{cacheKeyHash: []byte{1}, filePath: "file.spec.ts", originalName: "test1", suiteHierarchy: "suite"},
		}
		cached := map[string]*entity.CacheEntry{
			hex.EncodeToString([]byte{1}): {ConvertedName: "cached"},
		}

		result := groupUncachedByFile(metas, cached)

		if len(result) != 0 {
			t.Errorf("expected empty result when all cached, got %d files", len(result))
		}
	})

	t.Run("suites are properly grouped per file", func(t *testing.T) {
		metas := []testMeta{
			{cacheKeyHash: []byte{1}, filePath: "file.spec.ts", originalName: "test1", suiteHierarchy: "suiteA"},
			{cacheKeyHash: []byte{2}, filePath: "file.spec.ts", originalName: "test2", suiteHierarchy: "suiteB"},
			{cacheKeyHash: []byte{3}, filePath: "file.spec.ts", originalName: "test3", suiteHierarchy: "suiteA"},
		}
		cached := make(map[string]*entity.CacheEntry)

		result := groupUncachedByFile(metas, cached)

		file := result["file.spec.ts"]
		if file == nil {
			t.Fatal("file.spec.ts not found")
		}
		if len(file.suites) != 2 {
			t.Errorf("expected 2 suites, got %d", len(file.suites))
		}
	})

	_ = codebaseID
}
