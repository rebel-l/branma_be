package repositorystore_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/bootstrap"
	"github.com/rebel-l/branma_be/repository/repositorystore"
	"github.com/rebel-l/go-utils/osutils"
)

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 0. setup
	storagePath := filepath.Join(".", "..", "..", "storage", "test_repository")
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

	defer func() {
		if err = db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	// 3. test
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.actual.Create(context.Background(), db)
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("expected error '%v' but got '%v'", testCase.expectedErr, err)
			} else {
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

					fmt.Println(testCase.actual)
				}
			}
		})
	}
}
