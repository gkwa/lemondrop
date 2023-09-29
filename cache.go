package lemondrop

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
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
