package storage

import (
	"os"
	"path/filepath"
	"yoru/shared"
)

func getDatabasePath() (string, error) {
	baseDir, err := getBaseDirectory()
	if err != nil {
		return "", err
	}

	databaseName := shared.PackageName + ".db"

	return filepath.Join(baseDir, databaseName), ensureDirectoryExists(baseDir)
}

func getBaseDirectory() (string, error) {
	if shared.Version == "dev" {
		return os.Getwd()
	}
	return os.UserConfigDir()
}

func ensureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}
