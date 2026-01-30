package gemini

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"sync"
	"time"

	"github.com/specvital/worker/internal/domain/specview"
)

const (
	// taxonomyCacheTTL is the time-to-live for cached taxonomy entries.
	// 1 hour balances cache benefits vs stale data risk.
	taxonomyCacheTTL = 1 * time.Hour

	// taxonomyCacheMaxSize is the maximum number of entries in the cache.
	// Prevents unbounded memory growth.
	taxonomyCacheMaxSize = 1000
)

// TaxonomyCacheKey identifies a unique taxonomy extraction result.
// Cache key is generated from: analysis_id + language + model_id
type TaxonomyCacheKey struct {
	AnalysisID string
	Language   specview.Language
	ModelID    string
}

// taxonomyCacheEntry holds a cached taxonomy with expiration.
type taxonomyCacheEntry struct {
	expiresAt time.Time
	taxonomy  *specview.TaxonomyOutput
}

// TaxonomyCache provides in-memory caching for taxonomy extraction results.
// Thread-safe with TTL-based expiration.
type TaxonomyCache struct {
	entries map[string]taxonomyCacheEntry
	mu      sync.RWMutex
	now     func() time.Time // for testing
}

// NewTaxonomyCache creates a new taxonomy cache.
func NewTaxonomyCache() *TaxonomyCache {
	return &TaxonomyCache{
		entries: make(map[string]taxonomyCacheEntry),
		now:     time.Now,
	}
}

// Get retrieves a cached taxonomy if it exists and hasn't expired.
// Returns nil if not found or expired.
func (c *TaxonomyCache) Get(key TaxonomyCacheKey) *specview.TaxonomyOutput {
	hashKey := c.hashKey(key)

	c.mu.RLock()
	entry, exists := c.entries[hashKey]
	c.mu.RUnlock()

	if !exists {
		return nil
	}

	if c.now().After(entry.expiresAt) {
		// Expired - remove async to avoid blocking read path
		go c.deleteIfExpired(hashKey)
		return nil
	}

	slog.Debug("taxonomy cache hit",
		"analysis_id", key.AnalysisID,
		"language", key.Language,
		"model_id", key.ModelID,
	)

	return entry.taxonomy
}

// Set stores a taxonomy in the cache with TTL.
func (c *TaxonomyCache) Set(key TaxonomyCacheKey, taxonomy *specview.TaxonomyOutput) {
	if taxonomy == nil {
		return
	}

	hashKey := c.hashKey(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Enforce max size by evicting expired entries first
	if len(c.entries) >= taxonomyCacheMaxSize {
		c.evictExpired()
	}

	// If still at max, evict oldest entry
	if len(c.entries) >= taxonomyCacheMaxSize {
		c.evictOldest()
	}

	c.entries[hashKey] = taxonomyCacheEntry{
		expiresAt: c.now().Add(taxonomyCacheTTL),
		taxonomy:  taxonomy,
	}

	slog.Debug("taxonomy cached",
		"analysis_id", key.AnalysisID,
		"language", key.Language,
		"model_id", key.ModelID,
		"expires_at", c.entries[hashKey].expiresAt,
	)
}

// deleteIfExpired removes an entry only if it's still expired.
// Prevents race condition where a new valid entry was set after expiration check.
func (c *TaxonomyCache) deleteIfExpired(hashKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, exists := c.entries[hashKey]; exists && c.now().After(entry.expiresAt) {
		delete(c.entries, hashKey)
	}
}

// evictExpired removes all expired entries.
// Must be called with write lock held.
func (c *TaxonomyCache) evictExpired() {
	now := c.now()
	for k, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, k)
		}
	}
}

// evictOldest removes the entry with earliest expiration.
// Must be called with write lock held.
func (c *TaxonomyCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for k, entry := range c.entries {
		if oldestKey == "" || entry.expiresAt.Before(oldestTime) {
			oldestKey = k
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// hashKey generates a deterministic hash from cache key components.
func (c *TaxonomyCache) hashKey(key TaxonomyCacheKey) string {
	h := sha256.New()
	h.Write([]byte(key.AnalysisID))
	h.Write([]byte("|"))
	h.Write([]byte(key.Language))
	h.Write([]byte("|"))
	h.Write([]byte(key.ModelID))
	return hex.EncodeToString(h.Sum(nil))
}

// Size returns the current number of entries in the cache.
// Useful for testing and monitoring.
func (c *TaxonomyCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
