package repositorymodel_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/rebel-l/branma_be/repository/repositorymodel"
)

func TestRepository_DecodeJSON(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339Nano, "2019-12-31T03:36:57.9167778+01:00")
	modifiedAt, _ := time.Parse(time.RFC3339Nano, "2020-01-01T15:44:57.9168378+01:00")

	testCases := []struct {
		name        string
		actual      *repositorymodel.Repository
		json        io.Reader
		expected    *repositorymodel.Repository
		expectedErr error
	}{
		{
			name: "model is nil",
		},
		{
			name:        "no JSON format",
			actual:      &repositorymodel.Repository{},
			json:        bytes.NewReader([]byte("no JSON")),
			expected:    &repositorymodel.Repository{},
			expectedErr: repositorymodel.ErrDecodeJSON,
		},
		{
			name:   "success",
			actual: &repositorymodel.Repository{},
			json: bytes.NewReader([]byte(`{
				"id": 1,
				"name": "test",
				"url": "url",
				"created_at": "2019-12-31T03:36:57.9167778+01:00",
				"modified_at": "2020-01-01T15:44:57.9168378+01:00"
			}`)),
			expected: &repositorymodel.Repository{
				ID:         1,
				Name:       "test",
				URL:        "url",
				CreatedAt:  createdAt,
				ModifiedAt: modifiedAt,
			},
		},
		{
			name:     "empty json",
			actual:   &repositorymodel.Repository{},
			json:     bytes.NewReader([]byte("{}")),
			expected: &repositorymodel.Repository{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.actual.DecodeJSON(testCase.json)
			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)
				return
			}

			testRepository(t, testCase.expected, testCase.actual)
		})
	}
}

func testRepository(t *testing.T, expected, actual *repositorymodel.Repository) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected != nil && actual == nil || expected == nil && actual != nil {
		t.Errorf("expected repository to be '%v' but got '%v'", expected, actual)
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

	if !expected.CreatedAt.Equal(actual.CreatedAt) {
		t.Errorf("expected created at '%s' but got '%s'", expected.CreatedAt.String(), actual.CreatedAt.String())
	}

	if !expected.ModifiedAt.Equal(actual.ModifiedAt) {
		t.Errorf("expected modified at '%s' but got '%s'", expected.ModifiedAt.String(), actual.ModifiedAt.String())
	}
}
