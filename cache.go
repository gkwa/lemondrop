package lemondrop

import (
	"encoding/gob"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	gocache "github.com/patrickmn/go-cache"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
	"github.com/taylormonacelli/somespider"
)

var (
	regionsCache      *gocache.Cache
	cacheKey          string
	cacheRelativePath string
	expiration        time.Duration
	cleanupInterval   time.Duration
)

func init() {
	cacheKey = "aws-regions"
	cleanupInterval = 24 * time.Hour
	expiration = 12 * time.Hour
	regionsCache = gocache.New(expiration, cleanupInterval)
	cacheRelativePath = filepath.Join("lemondrop", "aws-regions.gob")
}

func fetchRegionsFromCache() (RegionDetails, error) {
	cachePath, err := somespider.GenPath(cacheRelativePath)
	if err != nil {
		return RegionDetails{}, err
	}

	if !mymazda.FileExists(cachePath) {
		return RegionDetails{}, nil
	}

	// unmarshal cache from file:
	file, err := os.Open(cachePath)
	if err != nil {
		slog.Debug("file access", "error", err.Error())
		return RegionDetails{}, err
	}
	defer file.Close()

	gobDecoder := gob.NewDecoder(file)
	gob.Register(RegionDetails{})

	var cacheMap map[string]gocache.Item
	if err := gobDecoder.Decode(&cacheMap); err != nil {
		slog.Debug("decode", "error", err.Error())
		return RegionDetails{}, err
	}

	cache2 := gocache.NewFrom(expiration, cleanupInterval, cacheMap)
	reply, future, found := cache2.GetWithExpiration(cacheKey)

	if reply == nil {
		// cache expired
		return RegionDetails{}, nil
	}

	expires := time.Until(future).Truncate(time.Second)
	e := reply.(RegionDetails)
	slog.Debug("newCache", "found", found, "expires", expires, "now", time.Now(), "future", future, "result", e)

	return e, nil
}
