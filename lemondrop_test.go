package lemondrop

import (
	"testing"
	"time"
)

func TestReadRegionsFromCache_Expired(t *testing.T) {
	// Set the mock timestamp provider with a time that exceeds the cache expiration
	mockTime := time.Now().Add(-48 * time.Hour)
	mockProvider := MockTimestampProvider{MockTime: mockTime}

	regions, err := readRegionsFromCache(mockProvider)
	if err == nil {
		t.Error("Expected error due to expired cache, but got no error")
	}

	if regions != nil {
		t.Error("Expected nil regions due to expired cache, but got non-nil regions")
	}
}
