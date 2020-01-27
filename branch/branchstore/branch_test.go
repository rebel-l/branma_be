package branchstore_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/branch/branchstore"
	"github.com/rebel-l/branma_be/repository/repositorystore"
	"github.com/rebel-l/branma_be/test"
)

const (
	testCluster = "test_branch"
)

func TestBranch_Create(t *testing.T) { // nolint:funlen
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := test.Setup(t, testCluster, "storeCreate")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	repo := &repositorystore.Repository{Name: "testrepo", URL: "testrepo.url"}
	if err := repo.Create(context.Background(), db); err != nil {
		t.Fatalf("preparing data failed: %v", err)
	}

	// 2. test
	testCases := []struct {
		name        string
		actual      *branchstore.Branch
		expected    *branchstore.Branch
		expectedErr error
	}{
		{
			name:        "branch is nil",
			expectedErr: branchstore.ErrDataMissing,
		},
		{
			name:        "branch has no name",
			actual:      &branchstore.Branch{TicketID: "JIRA-1", RepositoryID: 1},
			expectedErr: branchstore.ErrDataMissing,
		},
		{
			name:        "branch has no repository id",
			actual:      &branchstore.Branch{Name: "JIRA-1"},
			expectedErr: branchstore.ErrDataMissing,
		},
		{
			name:        "branch has ID",
			actual:      &branchstore.Branch{ID: 1, Name: "myname", RepositoryID: 1},
			expectedErr: branchstore.ErrIDIsSet,
		},
		{
			name: "success",
			actual: &branchstore.Branch{
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expected: &branchstore.Branch{
				ID:             1,
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
		},
		{
			name: "duplicate",
			actual: &branchstore.Branch{
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expectedErr: errors.New("UNIQUE constraint failed: branches.branch_name, branches.repository_id"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.actual.Create(context.Background(), db)
			test.CheckErrors(t, testCase.expectedErr, err)
			testBranch(t, testCase.expected, testCase.actual)
		})
	}
}

func TestBranch_Read(t *testing.T) { // nolint:funlen
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := test.Setup(t, testCluster, "storeRead")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	repo := &repositorystore.Repository{Name: "testreporead", URL: "testreporead.url"}
	if err := repo.Create(context.Background(), db); err != nil {
		t.Fatalf("preparing data failed: %v", err)
	}

	// 2. test
	testCases := []struct {
		name        string
		prepare     *branchstore.Branch
		actual      *branchstore.Branch
		expected    *branchstore.Branch
		expectedErr error
	}{
		{
			name:        "branch is nil",
			expectedErr: repositorystore.ErrIDMissing,
		},
		{
			name:        "ID not set",
			expectedErr: branchstore.ErrIDMissing,
			actual:      &branchstore.Branch{},
		},
		{
			name: "success",
			prepare: &branchstore.Branch{
				Name:           "mybranch",
				TicketID:       "ID-123",
				ParentTicketID: "EPIC-123",
				RepositoryID:   1,
				TicketSummary:  "create a branch manager",
				TicketStatus:   "done",
				TicketType:     "story",
			},
			actual: &branchstore.Branch{ID: 1},
			expected: &branchstore.Branch{
				ID:             1,
				Name:           "mybranch",
				TicketID:       "ID-123",
				ParentTicketID: "EPIC-123",
				RepositoryID:   1,
				TicketSummary:  "create a branch manager",
				TicketStatus:   "done",
				TicketType:     "story",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.prepare != nil {
				_ = testCase.prepare.Create(context.Background(), db)
			}

			err := testCase.actual.Read(context.Background(), db)
			test.CheckErrors(t, testCase.expectedErr, err)
			testBranch(t, testCase.expected, testCase.actual)
		})
	}
}

func TestBranch_Update(t *testing.T) { // nolint:funlen
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := test.Setup(t, testCluster, "storeUpdate")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	repo := &repositorystore.Repository{Name: "testrepoupdate", URL: "testrepoupdate.url"}
	if err := repo.Create(context.Background(), db); err != nil {
		t.Fatalf("preparing data failed: %v", err)
	}

	// 2. test
	testCases := []struct {
		name        string
		prepare     *branchstore.Branch
		actual      *branchstore.Branch
		expected    *branchstore.Branch
		expectedErr error
	}{
		{
			name:        "branch is nil",
			expectedErr: branchstore.ErrDataMissing,
		},
		{
			name: "branch has no name",
			actual: &branchstore.Branch{
				ID:             1,
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expectedErr: branchstore.ErrDataMissing,
		},
		{
			name: "branch has no repository id",
			actual: &branchstore.Branch{
				ID:             1,
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expectedErr: repositorystore.ErrDataMissing,
		},
		{
			name: "branch has no ID",
			actual: &branchstore.Branch{
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expectedErr: branchstore.ErrIDMissing,
		},
		{
			name: "success",
			prepare: &branchstore.Branch{
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			actual: &branchstore.Branch{
				ID:             1,
				Name:           "new name",
				TicketID:       "NEW-1",
				ParentTicketID: "NEW-2",
				RepositoryID:   1,
				TicketSummary:  "a new summary",
				TicketStatus:   "done",
				TicketType:     "bug",
				Closed:         false,
			},
			expected: &branchstore.Branch{
				ID:             1,
				Name:           "new name",
				TicketID:       "NEW-1",
				ParentTicketID: "NEW-2",
				RepositoryID:   1,
				TicketSummary:  "a new summary",
				TicketStatus:   "done",
				TicketType:     "bug",
				Closed:         false,
			},
		},
		{
			name: "not existing branch",
			actual: &branchstore.Branch{
				ID:             2,
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
				Closed:         true,
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.prepare != nil {
				_ = testCase.prepare.Create(context.Background(), db)
				time.Sleep(1 * time.Second)
			}

			err := testCase.actual.Update(context.Background(), db)
			test.CheckErrors(t, testCase.expectedErr, err)
			testBranch(t, testCase.expected, testCase.actual)

			if testCase.prepare != nil && testCase.actual != nil {
				if testCase.prepare.CreatedAt != testCase.actual.CreatedAt {
					t.Errorf(
						"expected created at '%s' but got '%s'",
						testCase.prepare.CreatedAt.String(),
						testCase.actual.CreatedAt.String(),
					)
				}

				if testCase.prepare.ModifiedAt.After(testCase.actual.ModifiedAt) {
					t.Errorf(
						"expected modified at '%s' to be before but got '%s'",
						testCase.prepare.ModifiedAt.String(),
						testCase.actual.ModifiedAt.String(),
					)
				}
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) { // nolint:funlen
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := test.Setup(t, testCluster, "storeDelete")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	repo := &repositorystore.Repository{Name: "testrepodelete", URL: "testrepodelete.url"}
	if err := repo.Create(context.Background(), db); err != nil {
		t.Fatalf("preparing data failed: %v", err)
	}

	// 2. test
	testCases := []struct {
		name        string
		prepare     *branchstore.Branch
		actual      *branchstore.Branch
		expectedErr error
	}{
		{
			name:        "branch is nil",
			expectedErr: branchstore.ErrIDMissing,
		},
		{
			name:        "branch has no ID",
			actual:      &branchstore.Branch{},
			expectedErr: branchstore.ErrIDMissing,
		},
		{
			name: "success",
			prepare: &branchstore.Branch{
				Name:           "myname",
				TicketID:       "JIRA-1",
				ParentTicketID: "JIRA-2",
				RepositoryID:   1,
				TicketSummary:  "a nice summary",
				TicketStatus:   "in progress",
				TicketType:     "improvement",
			},
			actual: &branchstore.Branch{ID: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var id int
			if testCase.prepare != nil {
				err := testCase.prepare.Create(context.Background(), db)
				if err != nil {
					t.Errorf("preparation failed: %v", err)
					return
				}
				id = testCase.prepare.ID
			}

			err := testCase.actual.Delete(context.Background(), db)
			test.CheckErrors(t, testCase.expectedErr, err)

			if id > 0 {
				testCase.actual = &branchstore.Branch{ID: id}
				err := testCase.actual.Read(context.Background(), db)
				if !errors.Is(err, sql.ErrNoRows) {
					t.Errorf("expected error '%v' after deletion but got '%v'", sql.ErrNoRows, err)
				}
			}
		})
	}
}

func testBranch(t *testing.T, expected, actual *branchstore.Branch) {
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
		t.Errorf("expected name '%s' but got '%s'", expected.Name, actual.Name)
	}

	if expected.TicketID != actual.TicketID {
		t.Errorf("expected ticket ID '%s' but got '%s'", expected.TicketID, actual.TicketID)
	}

	if expected.ParentTicketID != actual.ParentTicketID {
		t.Errorf("expected parent ticket ID '%s' but got '%s'", expected.ParentTicketID, actual.ParentTicketID)
	}

	if expected.RepositoryID != actual.RepositoryID {
		t.Errorf("expected repository ID '%d' but got '%d'", expected.RepositoryID, actual.RepositoryID)
	}

	if expected.TicketSummary != actual.TicketSummary {
		t.Errorf("expected ticket summary '%s' but got '%s'", expected.TicketSummary, actual.TicketSummary)
	}

	if expected.TicketStatus != actual.TicketStatus {
		t.Errorf("expected ticket status '%s' but got '%s'", expected.TicketStatus, actual.TicketStatus)
	}

	if expected.TicketType != actual.TicketType {
		t.Errorf("expected ticket type '%s' but got '%s'", expected.TicketType, actual.TicketType)
	}

	if expected.Closed != actual.Closed {
		t.Errorf("expected close '%t' but got '%t'", expected.Closed, actual.Closed)
	}

	if actual.CreatedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}

	if actual.ModifiedAt.IsZero() {
		t.Error("modified at should be greater than the zero date")
	}
}
