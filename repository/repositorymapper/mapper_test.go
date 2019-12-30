package repositorymapper_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/repository/repositorymapper"

	"github.com/rebel-l/branma_be/repository/repositorymodel"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/bootstrap"
	"github.com/rebel-l/go-utils/osutils"
)

func setup(t *testing.T, name string) *sqlx.DB {
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

	return db
}

func TestMapper_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "mapperLoad")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

}

func TestMapper_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "mapperSave")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	mapper := repositorymapper.New(db)

	// 2. test
	testCases := []struct {
		name        string
		actual      *repositorymodel.Repository
		expected    *repositorymodel.Repository
		expectedErr error
	}{
		{
			name:        "model is nil",
			expectedErr: repositorymapper.ErrNoData,
		},
		{
			name:     "model has no ID",
			actual:   &repositorymodel.Repository{Name: "myname", URL: "myurl"},
			expected: &repositorymodel.Repository{ID: 1, Name: "myname", URL: "myurl"},
		},
		{
			name:     "model has ID",
			actual:   &repositorymodel.Repository{ID: 1, Name: "newname", URL: "newurl"},
			expected: &repositorymodel.Repository{ID: 1, Name: "newname", URL: "newurl"},
		},
		{
			name:        "model is duplicate",
			actual:      &repositorymodel.Repository{Name: "newname", URL: "newurl"},
			expectedErr: repositorymapper.ErrSaveToDB,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := mapper.Save(context.Background(), testCase.actual)
			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)
			}

			testRepository(t, testCase.expected, res)
		})
	}
}

func testRepository(t *testing.T, expected, actual *repositorymodel.Repository) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected != nil && actual == nil || expected == nil && actual != nil {
		t.Errorf("expected repository '%v' but got '%v'", expected, actual)
		return
	}

	if expected.ID != actual.ID {
		t.Errorf("expected ID %d but got %d", expected.ID, actual.ID)
	}

	if expected.Name != actual.Name {
		t.Errorf("expected name '%s' but got '%s'", expected.Name, actual.Name)
	}

	if expected.URL != actual.URL {
		t.Errorf("expected URL '%s' but got '%s'", expected.URL, actual.URL)
	}

	if actual.CreatedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}

	if actual.ModifiedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}
}
