package lemondrop

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
)

var (
	regionsCache *cache.Cache
	cacheKey     string
)

func init() {
	regionsCache = cache.New(24*time.Hour, 24*time.Hour)
	cacheKey = "aws/regions"
}

func fetchFromCache() (RegionDetails, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return RegionDetails{}, err
	}

	regionsCache.LoadFile(cachePath)

	regions := make(RegionDetails)

	regionInterface, found := regionsCache.Get(cacheKey)
	if !found {
		return RegionDetails{}, nil
	}

	jsonData := regionInterface.(string)
	err = json.Unmarshal([]byte(jsonData), &regions)
	if err != nil {
		return RegionDetails{}, err
	}
	return regions, nil
}

func GetRegionDetails() (RegionDetails, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return RegionDetails{}, err
	}

	regions, err := fetchFromCache()
	if err != nil {
		return RegionDetails{}, err
	}

	if len(regions) != 0 {
		slog.Info("cache hit")
		return regions, nil
	}

	slog.Info("cache miss")

	regions, err = GetAllAwsRegions()
	if err != nil {
		return RegionDetails{}, err
	}

	jsonBytes, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		return RegionDetails{}, err
	}
	regionsCache.Set(cacheKey, string(jsonBytes), cache.DefaultExpiration)
	regionsCache.SaveFile(cachePath)

	return regions, nil
}
