package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rebel-l/smis"

	"github.com/rebel-l/branma_be/repository/repositorymapper"
	"github.com/rebel-l/branma_be/repository/repositorymodel"
)

type tcGet struct {
	name            string
	request         *http.Request
	expectedCode    int
	expectedPayload *Payload
}

func getTestCasesGet(t *testing.T) []tcGet { // nolint: funlen
	t.Helper()

	var testCases []tcGet

	// 1.
	req, err := http.NewRequest(http.MethodGet, "/repository/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := tcGet{
		name:         "success",
		request:      req,
		expectedCode: http.StatusOK,
		expectedPayload: &Payload{
			Repository: &repositorymodel.Repository{
				ID:   1,
				Name: "repo",
				URL:  "url",
			},
		},
	}

	testCases = append(testCases, c)

	// 2.
	req, err = http.NewRequest(http.MethodGet, "/repository/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tcGet{
		name:            "repository not found",
		request:         req,
		expectedCode:    http.StatusNotFound,
		expectedPayload: &Payload{Error: "failed to load repository for id: 3"},
	}

	testCases = append(testCases, c)

	// 3.
	req, err = http.NewRequest(http.MethodGet, "/repository/abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tcGet{
		name:            "id not integer",
		request:         req,
		expectedCode:    http.StatusBadRequest,
		expectedPayload: &Payload{Error: "converting id to integer failed"},
	}

	testCases = append(testCases, c)

	return testCases
}

func TestHandler_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	svc, db := setup(t, "endpointGet")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	if err := Init(svc, db); err != nil {
		t.Fatalf("failed to init routes: %v", err)
	}

	// 2. prepare test data
	mapper := repositorymapper.New(db)
	if _, err := mapper.Save(context.Background(), &repositorymodel.Repository{Name: "repo", URL: "url"}); err != nil {
		t.Fatalf("failed to prepare test data: %v", err)
	}

	// 3. test
	for _, testCase := range getTestCasesGet(t) {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			svc.Router.ServeHTTP(w, testCase.request)

			if testCase.expectedCode != w.Code {
				t.Errorf("expected code %d but got %d", testCase.expectedCode, w.Code)
			}

			contentType := w.Header().Get(smis.HeaderKeyContentType)
			if contentType != smis.HeaderContentTypeJSON {
				t.Errorf("expected content type '%s' but got '%s'", smis.HeaderContentTypeJSON, contentType)
			}

			actual := &Payload{}
			if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
				t.Fatalf("failed to decode json: %v | %s", err, w.Body.Bytes())
			}

			testPayload(t, testCase.expectedPayload, actual)
		})
	}
}

func TestHandler_Get_RequestNil(t *testing.T) {
	handler := Handler{}
	w := httptest.NewRecorder()
	handler.get(w, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected code %d but got %d", http.StatusBadRequest, w.Code)
	}

	actual := &Payload{}
	if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
		t.Fatalf("failed to decode json: %v | %s", err, w.Body.Bytes())
	}

	testPayload(t, &Payload{Error: "request is empty"}, actual)
}

func TestHandler_Get_NoID(t *testing.T) {
	handler := New(&smis.Service{}, nil)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/repository/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected code %d but got %d", http.StatusBadRequest, w.Code)
	}

	actual := &Payload{}
	if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
		t.Fatalf("failed to decode json: %v | %s", err, w.Body.Bytes())
	}

	testPayload(t, &Payload{Error: "request is empty"}, actual)
}
