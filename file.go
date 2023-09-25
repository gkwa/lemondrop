package lemondrop

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func getCachePath() (string, error) {
	configRelPath := "lemondrop/regions.db"
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}

	dirPerm := os.FileMode(0700)

	d := filepath.Dir(configFilePath)

	if err := os.MkdirAll(d, dirPerm); err != nil {
		return "", err
	}

	return configFilePath, nil
}
