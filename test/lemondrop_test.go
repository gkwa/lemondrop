package lemondrop_test

import (
	"os"
	"testing"
	"time"
)

func TestCacheExpiration(t *testing.T) {
	// Simulate an expired cache by setting a timestamp in the past
	expiredTimestamp := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)

	// Modify the cacheFilePath to a temporary file for testing
	lemondrop.CacheFilePath = "temporary_cache.json"

	// Write a cache file with an expired timestamp
	err := lemondrop.WriteRegionsToCacheWithTimestamp([]lemondrop.Region{}, expiredTimestamp)
	if err != nil {
		t.Fatalf("Failed to write cache file: %s", err)
	}

	// Test reading the regions from cache
	_, err = lemondrop.ReadRegionsFromCache()
	if err == nil {
		t.Error("Expected cache expiration error, but got no error")
	} else if err.Error() != "cache has expired" {
		t.Errorf("Expected 'cache has expired' error, but got '%s'", err.Error())
	}

	// Clean up the temporary cache file
	err = os.Remove("temporary_cache.json")
	if err != nil {
		t.Errorf("Failed to clean up cache file: %s", err)
	}
}
