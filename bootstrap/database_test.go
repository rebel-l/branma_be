package bootstrap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/config"

	"github.com/rebel-l/go-utils/slice"

	"github.com/rebel-l/go-utils/osutils"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/bootstrap"
)

func setup(t *testing.T, name string) *config.Database {
	t.Helper()

	// 1. config
	storagePath := filepath.Join(".", "..", "storage", name)
	scriptPath := filepath.Join(".", "..", "scripts", "schema")
	conf := &config.Database{
		StoragePath:       &storagePath,
		SchemaScriptsPath: &scriptPath,
	}

	// 2. clean up
	if osutils.FileOrPathExists(conf.GetStoragePath()) {
		if err := os.RemoveAll(conf.GetStoragePath()); err != nil {
			t.Fatalf("failed to cleanup test files: %v", err)
		}
	}

	return conf
}

func TestDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	fixtures := slice.StringSlice{
		"schema_script",
		"sqlite_sequence",
		"branches",
		"versions",
		"branch_versions",
		"commits",
		"branch_commits",
		"repositories",
	}

	// 1. setup
	conf := setup(t, "test_bootstrap")

	// 2. do the test
	db, err := bootstrap.Database(conf, "0.0.0")
	if err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	// 3. do the assertions
	var tables slice.StringSlice

	q := db.Rebind("SELECT name FROM sqlite_master WHERE type='table';")

	if err = db.Select(&tables, q); err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	if !fixtures.IsEqual(tables) {
		t.Errorf("tables are not created, expected: '%v' | got: '%v'", fixtures, tables)
	}
}

func TestDatabaseReset(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	fixtures := slice.StringSlice{
		"schema_script",
		"sqlite_sequence",
	}

	// 1. setup
	conf := setup(t, "test_reset")

	// 2. do the test
	db, err := bootstrap.Database(conf, "0.0.0")
	if err != nil {
		t.Fatalf("No error expected on bbotstrap: %v", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	if err = bootstrap.DatabaseReset(conf); err != nil {
		t.Fatalf("No error expected on reset: %v", err)
	}

	// 3. do the assertions
	var tables slice.StringSlice

	q := db.Rebind("SELECT name FROM sqlite_master WHERE type='table';")

	if err = db.Select(&tables, q); err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	if !fixtures.IsEqual(tables) {
		t.Errorf("tables are not reseted, expected: '%v' | got: '%v'", fixtures, tables)
	}
}
