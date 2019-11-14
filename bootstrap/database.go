package bootstrap

import (
	"path/filepath"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	storageFileName = "branma.db"
)

// Database initialises the database and returns the connection
func Database(storagePath, scriptsPath string) (*sqlx.DB, error) {
	fileName, err := createStorage(storagePath)
	if err != nil {
		return nil, err // TODO: pimp error
	}

	db, err := sqlx.Open("sqlite3", fileName)
	if err != nil {
		return nil, err // TODO: pimp error
	}

	return db, nil
}

func createStorage(path string) (string, error) {
	if err := osutils.CreateDirectoryIfNotExists(path); err != nil {
		return "", err // TODO: pimp error
	}

	fileName := filepath.Join(path, storageFileName)
	if err := osutils.CreateFileIfNotExists(fileName); err != nil {
		return "", err // TODO: pimp error
	}

	return fileName, nil
}
