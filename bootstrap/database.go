package bootstrap

import (
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
		return nil, err // TODO: pimp error
	}

	db, err := sqlx.Open("sqlite3", fileName)
	if err != nil {
		return nil, err // TODO: pimp error
	}

	err = createSchema(db, scriptPath, version)

	return db, err // TODO: pimp error
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

func createSchema(db *sqlx.DB, scriptPath, version string) error {
	s := schema.New(db)
	s.WithProgressBar()
	return s.Upgrade(scriptPath, version) // TODO: check error and revert last
}
