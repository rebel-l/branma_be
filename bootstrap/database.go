package bootstrap

import (
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	storageFileName = "branma.db"
)

// Database initialises the database and returns the connection
func Database(storagePath, scriptsPath string) error {
	if err := createStorage(storagePath); err != nil {
		return err // TODO: pimp error
	}

	return nil
}

func createStorage(path string) error {
	if err := osutils.CreateDirectoryIfNotExists(path); err != nil {
		return err // TODO: pimp error
	}

	fileName := filepath.Join(path, storageFileName)
	if err := osutils.CreateFileIfNotExists(fileName); err != nil {
		return err // TODO: pimp error
	}

	return nil
}
