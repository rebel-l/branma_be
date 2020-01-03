package repository

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rebel-l/smis"

	"github.com/rebel-l/branma_be/repository/repositorymodel"

	_ "github.com/mattn/go-sqlite3"
)

type tcPut struct {
	name            string
	request         *http.Request
	expectedPayload *Payload
	expectedStatus  int
}

func getTestCasesPut(t *testing.T) []tcPut { // nolint:funlen
	t.Helper()

	var testCases []tcPut

	// 1.
	c := tcPut{
		name:            "request nil",
		request:         nil,
		expectedStatus:  http.StatusBadRequest,
		expectedPayload: &Payload{Error: "request is empty"},
	}
	testCases = append(testCases, c)

	// 2.
	req, err := http.NewRequest(http.MethodPut, "/repository", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tcPut{
		name:            "request body nil",
		request:         req,
		expectedStatus:  http.StatusBadRequest,
		expectedPayload: &Payload{Error: "request body is empty"},
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

	c = tcPut{
		name:    "new repository",
		request: req,
		expectedPayload: NewPayload(&repositorymodel.Repository{
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

	c = tcPut{
		name:    "update repository",
		request: req,
		expectedPayload: NewPayload(&repositorymodel.Repository{
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

	ep := New(svc, db)
	handler := http.HandlerFunc(ep.put)

	// 2. test
	for _, testCase := range getTestCasesPut(t) {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, testCase.request)

			if testCase.expectedStatus != w.Code {
				t.Errorf("expected code %d but got %d", testCase.expectedStatus, w.Code)
			}

			contentType := w.Header().Get(smis.HeaderKeyContentType)
			if contentType != smis.HeaderContentTypeJSON {
				t.Errorf("expected content type '%s' but got '%s'", smis.HeaderContentTypeJSON, contentType)
			}

			actual := &Payload{}
			if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
				t.Fatalf("failed to decode json: %v", err)
			}

			testPayload(t, testCase.expectedPayload, actual)
		})
	}
}
