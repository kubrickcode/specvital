package gemini

import (
	"fmt"
	"testing"
	"time"

	"github.com/specvital/worker/internal/domain/specview"
)

func TestTaxonomyCache_Hit_SkipsAPI(t *testing.T) {
	cache := NewTaxonomyCache()

	key := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{
				Description: "Test domain",
				Name:        "Authentication",
				Features: []specview.TaxonomyFeature{
					{FileIndices: []int{0, 1}, Name: "Login"},
				},
			},
		},
	}

	// Set taxonomy in cache
	cache.Set(key, taxonomy)

	// Get should return cached value
	cached := cache.Get(key)
	if cached == nil {
		t.Fatal("expected cache hit, got nil")
	}

	if len(cached.Domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(cached.Domains))
	}

	if cached.Domains[0].Name != "Authentication" {
		t.Errorf("expected domain name 'Authentication', got %q", cached.Domains[0].Name)
	}
}

func TestTaxonomyCache_Miss_ReturnsNil(t *testing.T) {
	cache := NewTaxonomyCache()

	key := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	// Get without setting should return nil
	cached := cache.Get(key)
	if cached != nil {
		t.Errorf("expected cache miss (nil), got %v", cached)
	}
}

func TestTaxonomyCache_DifferentKey_CacheMiss(t *testing.T) {
	cache := NewTaxonomyCache()

	key1 := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	key2 := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "English", // Different language
		ModelID:    "gemini-2.5-flash",
	}

	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{Name: "Test", Features: []specview.TaxonomyFeature{}},
		},
	}

	cache.Set(key1, taxonomy)

	// Different language should miss
	cached := cache.Get(key2)
	if cached != nil {
		t.Errorf("expected cache miss for different language, got %v", cached)
	}

	// Original key should still hit
	cached = cache.Get(key1)
	if cached == nil {
		t.Error("expected cache hit for original key")
	}
}

func TestTaxonomyCache_TTLExpiration(t *testing.T) {
	cache := NewTaxonomyCache()

	// Mock time for testing
	currentTime := time.Now()
	cache.now = func() time.Time { return currentTime }

	key := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{Name: "Test", Features: []specview.TaxonomyFeature{}},
		},
	}

	cache.Set(key, taxonomy)

	// Should hit immediately
	if cache.Get(key) == nil {
		t.Error("expected cache hit before TTL")
	}

	// Advance time past TTL
	currentTime = currentTime.Add(taxonomyCacheTTL + time.Minute)

	// Should miss after TTL
	if cache.Get(key) != nil {
		t.Error("expected cache miss after TTL expiration")
	}
}

func TestTaxonomyCache_NilTaxonomy_NotCached(t *testing.T) {
	cache := NewTaxonomyCache()

	key := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	// Setting nil should be no-op
	cache.Set(key, nil)

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after setting nil, got %d", cache.Size())
	}
}

func TestTaxonomyCache_MaxSize_Eviction(t *testing.T) {
	cache := NewTaxonomyCache()

	// Mock time for deterministic ordering
	currentTime := time.Now()
	cache.now = func() time.Time { return currentTime }

	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{Name: "Test", Features: []specview.TaxonomyFeature{}},
		},
	}

	// Fill cache to max size
	for i := 0; i < taxonomyCacheMaxSize; i++ {
		key := TaxonomyCacheKey{
			AnalysisID: fmt.Sprintf("analysis-%d", i),
			Language:   "Korean",
			ModelID:    "gemini-2.5-flash",
		}
		cache.Set(key, taxonomy)
		// Advance time slightly for each entry to ensure ordering
		currentTime = currentTime.Add(time.Millisecond)
	}

	if cache.Size() != taxonomyCacheMaxSize {
		t.Errorf("expected cache size %d, got %d", taxonomyCacheMaxSize, cache.Size())
	}

	// Adding one more should trigger eviction
	newKey := TaxonomyCacheKey{
		AnalysisID: "new-analysis",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}
	cache.Set(newKey, taxonomy)

	// Size should not exceed max
	if cache.Size() > taxonomyCacheMaxSize {
		t.Errorf("cache size %d exceeds max %d", cache.Size(), taxonomyCacheMaxSize)
	}

	// New entry should be cached
	if cache.Get(newKey) == nil {
		t.Error("expected new entry to be cached")
	}
}

func TestTaxonomyCache_ConcurrentAccess(t *testing.T) {
	cache := NewTaxonomyCache()

	taxonomy := &specview.TaxonomyOutput{
		Domains: []specview.TaxonomyDomain{
			{Name: "Test", Features: []specview.TaxonomyFeature{}},
		},
	}

	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := TaxonomyCacheKey{
				AnalysisID: "analysis-concurrent",
				Language:   specview.Language("lang-" + string(rune('a'+idx))),
				ModelID:    "gemini-2.5-flash",
			}
			cache.Set(key, taxonomy)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := TaxonomyCacheKey{
				AnalysisID: "analysis-concurrent",
				Language:   specview.Language("lang-" + string(rune('a'+idx))),
				ModelID:    "gemini-2.5-flash",
			}
			_ = cache.Get(key)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Should not panic and cache should have entries
	if cache.Size() == 0 {
		t.Error("expected non-empty cache after concurrent access")
	}
}

func TestTaxonomyCache_HashKey_Deterministic(t *testing.T) {
	cache := NewTaxonomyCache()

	key := TaxonomyCacheKey{
		AnalysisID: "analysis-123",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	hash1 := cache.hashKey(key)
	hash2 := cache.hashKey(key)

	if hash1 != hash2 {
		t.Errorf("hash key not deterministic: %s != %s", hash1, hash2)
	}

	// Different key should produce different hash
	key2 := TaxonomyCacheKey{
		AnalysisID: "analysis-456",
		Language:   "Korean",
		ModelID:    "gemini-2.5-flash",
	}

	hash3 := cache.hashKey(key2)
	if hash1 == hash3 {
		t.Error("different keys should produce different hashes")
	}
}
