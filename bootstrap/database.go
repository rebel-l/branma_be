package bootstrap

import (
	"fmt"
	"path/filepath"

	"github.com/rebel-l/schema"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	storageFileName = "branma.db"
)

// Database initialises the database and returns the connection
func Database(storagePath, scriptPath, version string) (*sqlx.DB, error) {
	fileName, err := createStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("bootstrap database, create storage failed: %v", err)
	}

	db, err := sqlx.Open("sqlite3", fileName)
	if err != nil {
		return nil, fmt.Errorf("bootstrap database, open database failed: %w", err)
	}

	err = createSchema(db, scriptPath, version)

	return db, fmt.Errorf("bootstrap database, create schema failed: %w", err)
}

func createStorage(path string) (string, error) {
	if err := osutils.CreateDirectoryIfNotExists(path); err != nil {
		return "", err
	}

	fileName := filepath.Join(path, storageFileName)
	if err := osutils.CreateFileIfNotExists(fileName); err != nil {
		return "", err
	}

	return fileName, nil
}

func createSchema(db *sqlx.DB, scriptPath, version string) error {
	s := schema.New(db)
	s.WithProgressBar()

	return s.Upgrade(scriptPath, version)
}
