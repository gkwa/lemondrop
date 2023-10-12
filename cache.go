package lemondrop

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/taylormonacelli/somespider"
)

var (
	regionsCache *cache.Cache
	cacheKey     string
	relCachPath  string
)

func init() {
	regionsCache = cache.New(24*time.Hour, 24*time.Hour)
	cacheKey = "aws/regions"
	relCachPath = filepath.Join("lmondrop", "regions.db")
}

func fetchFromCache() (RegionDetails, error) {
	cachePath, err := somespider.GenPath(relCachPath)
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
