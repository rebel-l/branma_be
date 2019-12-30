package repository_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

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
	storagePath := filepath.Join(".", "..", "..", "storage", "test_repository", name)
	scriptPath := filepath.Join(".", "..", "..", "scripts", "schema")

	// 1. clean up
	if osutils.FileOrPathExists(storagePath) {
		if err := os.RemoveAll(storagePath); err != nil {
			t.Fatalf("failed to cleanup test files: %v", err)
		}
	}

	// 2. init database
	db, err := bootstrap.Database(storagePath, scriptPath, "0.0.0")
	if err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	svc, err := smis.NewService(&http.Server{}, mux.NewRouter(), logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	return svc, db
}
