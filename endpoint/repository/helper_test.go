package repository

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/config"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/rebel-l/branma_be/bootstrap"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/smis"
)

func setup(t *testing.T, name string) (*smis.Service, *sqlx.DB) {
	t.Helper()

	// 0. init path
	conf := &config.Database{
		StoragePath:       filepath.Join(".", "..", "..", "storage", "test_repository", name),
		SchemaScriptsPath: filepath.Join(".", "..", "..", "scripts", "schema"),
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

	svc, err := smis.NewService(&http.Server{}, mux.NewRouter(), logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	return svc, db
}

func testPayload(t *testing.T, expected, actual *Payload) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected != nil && actual == nil || expected == nil && actual != nil {
		t.Errorf("expected response to be '%v' but got '%v'", expected, actual)
		return
	}

	if expected.Repository == nil && actual.Repository == nil {
		return
	}

	if expected.Repository != nil && actual.Repository == nil ||
		expected.Repository == nil && actual.Repository != nil {
		t.Errorf("expected repository to be '%v' but got '%v'", expected.Repository, actual.Repository)
		return
	}

	if expected.Error != actual.Error {
		t.Errorf("expectedd error '%v' but got '%v'", expected.Error, actual.Error)
	}

	if expected.Repository.ID != actual.Repository.ID {
		t.Errorf("expected ID %d but got %d", expected.Repository.ID, actual.Repository.ID)
	}

	if expected.Repository.Name != actual.Repository.Name {
		t.Errorf("expected name %s but got %s", expected.Repository.Name, actual.Repository.Name)
	}

	if expected.Repository.URL != actual.Repository.URL {
		t.Errorf("expected URL %s but got %s", expected.Repository.URL, actual.Repository.URL)
	}

	if actual.Repository.CreatedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}

	if actual.Repository.ModifiedAt.IsZero() {
		t.Error("modified at should be greater than the zero date")
	}
}
