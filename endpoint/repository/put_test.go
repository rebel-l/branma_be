package repository_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rebel-l/branma_be/repository/repositorymodel"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/branma_be/endpoint/repository"
)

type tc struct {
	name            string
	request         *http.Request
	expectedPayload *repository.Payload
	expectedStatus  int
}

func getTestCases(t *testing.T) []tc { // nolint:funlen
	t.Helper()

	var testCases []tc

	// 1.
	c := tc{
		name:            "request nil",
		request:         nil,
		expectedStatus:  http.StatusBadRequest,
		expectedPayload: &repository.Payload{Error: "request is empty"},
	}
	testCases = append(testCases, c)

	// 2.
	req, err := http.NewRequest(http.MethodPut, "/repository", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tc{
		name:            "request body nil",
		request:         req,
		expectedStatus:  http.StatusBadRequest,
		expectedPayload: &repository.Payload{Error: "request body is empty"},
	}
	testCases = append(testCases, c)

	// 3.
	body := `{
		"name": "new",
		"url": "new url"
	}`

	req, err = http.NewRequest(http.MethodPut, "/repository", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	c = tc{
		name:    "new repository",
		request: req,
		expectedPayload: repository.NewPayload(&repositorymodel.Repository{
			ID:   1,
			Name: "new",
			URL:  "new url",
		}),
		expectedStatus: http.StatusCreated,
	}
	testCases = append(testCases, c)

	// 4.
	body = `{
		"id": 1,
		"name": "changed",
		"url": "changed url"
	}`

	req, err = http.NewRequest(http.MethodPut, "/repository", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	c = tc{
		name:    "update repository",
		request: req,
		expectedPayload: repository.NewPayload(&repositorymodel.Repository{
			ID:   1,
			Name: "changed",
			URL:  "changed url",
		}),
		expectedStatus: http.StatusOK,
	}
	testCases = append(testCases, c)

	return testCases
}

func Test_Put(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	svc, db := setup(t, "endpointPut")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	ep := repository.New(svc, db)
	handler := http.HandlerFunc(ep.Put)

	// 2. test
	for _, testCase := range getTestCases(t) {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, testCase.request)

			if testCase.expectedStatus != w.Code {
				t.Errorf("expected code %d but got %d", testCase.expectedStatus, w.Code)
			}

			actual := &repository.Payload{}
			if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
				t.Fatalf("failed to decode json: %v", err)
			}

			testPayload(t, testCase.expectedPayload, actual)
		})
	}
}

func testPayload(t *testing.T, expected, actual *repository.Payload) {
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
