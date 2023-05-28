package lemondrop

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestGetAllAwsRegions_CacheExpiration(t *testing.T) {
	// Mock the cache file with an expired timestamp
	cachedRegions := CachedRegions{
		Timestamp: time.Now().Add(-2 * cacheExpiration).Format(time.RFC3339),
		Regions:   []types.Region{},
	}
	err := writeRegionsToCache(cachedRegions.Regions)
	if err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	// Call the function under test
	regions, err := GetAllAwsRegions()
	if err != nil {
		t.Fatalf("Failed to get AWS regions: %v", err)
	}

	// Verify that the regions were fetched from the AWS API
	if len(regions) == 0 {
		t.Error("Expected non-empty list of regions")
	}

	// Verify that the cache file was regenerated
	_, err = readRegionsFromCache()
	if err != nil {
		t.Errorf("Failed to read cache file: %v", err)
	}
}
