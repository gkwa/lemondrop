package lemondrop

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/adrg/xdg"
)

func getCachePath() (string, error) {
	configRelPath := "lemondrop/regions.db"
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}

	dirPerm := os.FileMode(0o700)

	d := filepath.Dir(configFilePath)

	if err := os.MkdirAll(d, dirPerm); err != nil {
		return "", err
	}

	slog.Debug("cache", "path", configFilePath)
	logPathStats(configFilePath)
	return configFilePath, nil
}

func logPathStats(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("stat", "path", filePath, "error", err.Error())
		return
	}

	// Get the file's owner user ID
	fileUID := fileInfo.Sys().(*syscall.Stat_t).Uid

	// Use the user package to get the user information
	u, err := user.LookupId(fmt.Sprintf("%d", fileUID))
	if err != nil {
		slog.Error("user info", "user", u, "error", err.Error())
		return
	}

	slog.Debug("owner", "path", filePath, "user", u.Username)
}
