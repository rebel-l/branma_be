package repositorystore_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/bootstrap"
	"github.com/rebel-l/branma_be/repository/repositorystore"
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

func TestRepository_Create(t *testing.T) { // nolint:funlen
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "create")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	// 2. test
	testCases := []struct {
		name        string
		actual      *repositorystore.Repository
		expected    *repositorystore.Repository
		expectedErr error
	}{
		{
			name:        "repository is nil",
			expectedErr: repositorystore.ErrDataMissing,
		},
		{
			name:        "repository has no name",
			actual:      &repositorystore.Repository{URL: "myurl"},
			expectedErr: repositorystore.ErrDataMissing,
		},
		{
			name:        "repository has no url",
			actual:      &repositorystore.Repository{Name: "myname"},
			expectedErr: repositorystore.ErrDataMissing,
		},
		{
			name:        "repository has ID",
			actual:      &repositorystore.Repository{ID: 1, Name: "myname", URL: "myurl"},
			expectedErr: repositorystore.ErrIDIsSet,
		},
		{
			name:     "success",
			actual:   &repositorystore.Repository{Name: "myname", URL: "myurl"},
			expected: &repositorystore.Repository{Name: "myname", URL: "myurl"},
		},
		{
			name:        "duplicate",
			actual:      &repositorystore.Repository{Name: "myname", URL: "myurl"},
			expectedErr: errors.New("UNIQUE constraint failed: repositories.url, repositories.name"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.actual.Create(context.Background(), db)
			checkErrors(t, testCase.expectedErr, err)

			if testCase.expected != nil && testCase.actual != nil {
				if testCase.actual.ID <= 0 {
					t.Errorf("expected id set to value greater than 0 but got '%d'", testCase.actual.ID)
				}

				if testCase.actual.CreatedAt.IsZero() {
					t.Error("created at should be greater than the zero date")
				}

				if testCase.actual.ModifiedAt.IsZero() {
					t.Error("created at should be greater than the zero date")
				}

				if testCase.expected.Name != testCase.actual.Name {
					t.Errorf("expectade name '%s' but got '%s'", testCase.expected.Name, testCase.actual.Name)
				}

				if testCase.expected.URL != testCase.actual.URL {
					t.Errorf("expectade url '%s' but got '%s'", testCase.expected.URL, testCase.actual.URL)
				}
			}
		})
	}
}

func TestRepository_Read(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "read")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	// 2. test
	testCases := []struct {
		name        string
		prepare     *repositorystore.Repository
		actual      *repositorystore.Repository
		expected    *repositorystore.Repository
		expectedErr error
	}{
		{
			name:        "repository is nil",
			expectedErr: repositorystore.ErrIDMissing,
		},
		{
			name:        "ID not set",
			expectedErr: repositorystore.ErrIDMissing,
			actual:      &repositorystore.Repository{},
		},
		{
			name:     "success",
			prepare:  &repositorystore.Repository{Name: "project", URL: "myproject.git"},
			actual:   &repositorystore.Repository{ID: 1},
			expected: &repositorystore.Repository{ID: 1, Name: "project", URL: "myproject.git"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.prepare != nil {
				_ = testCase.prepare.Create(context.Background(), db)
			}

			err := testCase.actual.Read(context.Background(), db)
			checkErrors(t, testCase.expectedErr, err)
			testRepository(t, testCase.expected, testCase.actual)
		})
	}
}

func testRepository(t *testing.T, expected, actual *repositorystore.Repository) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected != nil && actual == nil || expected == nil && actual != nil {
		return
	}

	if expected.ID != actual.ID {
		t.Errorf("expected ID %d but got %d", expected.ID, actual.ID)
	}

	if expected.Name != actual.Name {
		t.Errorf("expectade name '%s' but got '%s'", expected.Name, actual.Name)
	}

	if expected.URL != actual.URL {
		t.Errorf("expectade url '%s' but got '%s'", expected.URL, actual.URL)
	}

	if actual.CreatedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}

	if actual.ModifiedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}
}

func TestRepository_Update(t *testing.T) {
	// TODO: implement
	t.Skip("not implemented")
}

func TestRepository_Delete(t *testing.T) {
	// TODO: implement
	t.Skip("not implemented")
}

func TestRepository_IsValid(t *testing.T) {
	// TODO: implement
	t.Skip("not implemented")
}

func checkErrors(t *testing.T, expected, actual error) {
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
