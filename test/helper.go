package test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/bootstrap"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/branma_be/config"
	"github.com/rebel-l/go-utils/osutils"
)

// Setup takes care about creating, cleaning and establishing database / connection
func Setup(t *testing.T, cluster, name string) *sqlx.DB {
	t.Helper()

	// 0. init path
	storagePath := filepath.Join(".", "..", "..", "storage", cluster, name)
	scriptPath := filepath.Join(".", "..", "..", "scripts", "schema")
	conf := &config.Database{
		StoragePath:       &storagePath,
		SchemaScriptsPath: &scriptPath,
	}

	// 1. clean up
	if osutils.FileOrPathExists(conf.GetStoragePath()) {
		if err := os.RemoveAll(conf.GetStoragePath()); err != nil {
			t.Fatalf("failed to cleanup test files: %v", err)
		}
	}

	// 2. init database
	db, err := bootstrap.Database(conf, "0.0.0")
	if err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	return db
}

// CheckErrors is a helper to assert errors
func CheckErrors(t *testing.T, expected, actual error) {
	t.Helper()

	if errors.Is(actual, expected) {
		return
	}

	if expected != nil && actual != nil {
		if expected.Error() != actual.Error() {
			t.Errorf("expected error '%v' but got '%v'", expected, actual)
		}

		return
	}

	t.Errorf("expected error '%v' but got '%v'", expected, actual)
}
