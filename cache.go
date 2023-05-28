package lemondrop

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	cacheFilePath   = "regions_cache.json"
	cacheExpiration = 1 * 24 * time.Hour
)

type CachedRegions struct {
	Timestamp string         `json:"timestamp"`
	Regions   []types.Region `json:"regions"`
}

type TimestampProvider interface {
	Now() time.Time
}

type DefaultTimestampProvider struct{}

func (p DefaultTimestampProvider) Now() time.Time {
	return time.Now()
}

type MockTimestampProvider struct {
	MockTime time.Time
}

func (p MockTimestampProvider) Now() time.Time {
	return p.MockTime
}

func readRegionsFromCache(timestampProvider TimestampProvider) ([]types.Region, error) {
	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer cacheFile.Close()

	var cachedRegions CachedRegions
	err = json.NewDecoder(cacheFile).Decode(&cachedRegions)
	if err != nil {
		return nil, err
	}

	cachedTime, err := time.Parse(time.RFC3339, cachedRegions.Timestamp)
	if err != nil {
		return nil, err
	}

	// Check if the cached regions have expired
	if timestampProvider.Now().Sub(cachedTime) > cacheExpiration {
		return nil, fmt.Errorf("cache has expired")
	}

	return cachedRegions.Regions, nil
}

func writeRegionsToCache(regions []types.Region) error {
	cacheFileDir := filepath.Dir(cacheFilePath)
	if err := os.MkdirAll(cacheFileDir, 0o755); err != nil {
		return err
	}

	cacheFile, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	cachedRegions := CachedRegions{
		Timestamp: time.Now().Format(time.RFC3339),
		Regions:   regions,
	}

	encodedRegions, err := json.MarshalIndent(cachedRegions, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFilePath, encodedRegions, 0o644)
	if err != nil {
		return err
	}

	return nil
}
