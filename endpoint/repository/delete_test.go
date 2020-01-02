package repository_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rebel-l/branma_be/repository/repositorymapper"
	"github.com/rebel-l/branma_be/repository/repositorymodel"

	"github.com/rebel-l/branma_be/endpoint/repository"
	"github.com/rebel-l/smis"
)

type tcDelete struct {
	name            string
	request         *http.Request
	expectedCode    int
	expectedPayload string
}

func getTestCasesDelete(t *testing.T) []tcDelete { // nolint: funlen
	t.Helper()

	var testCases []tcDelete

	// 1.
	req, err := http.NewRequest(http.MethodDelete, "/repository/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := tcDelete{
		name:            "success",
		request:         req,
		expectedCode:    http.StatusOK,
		expectedPayload: "{}",
	}

	testCases = append(testCases, c)

	// 2.
	req, err = http.NewRequest(http.MethodDelete, "/repository/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tcDelete{
		name:            "repository does not exist",
		request:         req,
		expectedCode:    http.StatusOK,
		expectedPayload: "{}",
	}

	testCases = append(testCases, c)

	// 3.
	req, err = http.NewRequest(http.MethodDelete, "/repository/abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	c = tcDelete{
		name:            "id not integer",
		request:         req,
		expectedCode:    http.StatusBadRequest,
		expectedPayload: `{"error":"converting id to integer failed"}`,
	}

	testCases = append(testCases, c)

	return testCases
}

func TestHandler_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	svc, db := setup(t, "endpointDelete")

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	}()

	if err := repository.Init(svc, db); err != nil {
		t.Fatalf("failed to init routes: %v", err)
	}

	// 2. prepare test data
	mapper := repositorymapper.New(db)
	if _, err := mapper.Save(context.Background(), &repositorymodel.Repository{Name: "repo", URL: "url"}); err != nil {
		t.Fatalf("failed to prepare test data: %v", err)
	}

	// 3. test
	for _, testCase := range getTestCasesDelete(t) {
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

			if testCase.expectedPayload != w.Body.String() {
				t.Errorf("expected payload %s but got %s", testCase.expectedPayload, w.Body.String())
			}
		})
	}
}

func TestHandler_Delete_RequestNil(t *testing.T) {
	handler := repository.Handler{}
	w := httptest.NewRecorder()
	handler.Delete(w, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected code %d but got %d", http.StatusBadRequest, w.Code)
	}

	actual := &repository.Payload{}
	if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
		t.Fatalf("failed to decode json: %v | %s", err, w.Body.Bytes())
	}

	testPayload(t, &repository.Payload{Error: "request is empty"}, actual)
}

func TestHandler_Delete_NoID(t *testing.T) {
	handler := repository.New(&smis.Service{}, nil)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodDelete, "/repository/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected code %d but got %d", http.StatusBadRequest, w.Code)
	}

	actual := &repository.Payload{}
	if err := json.Unmarshal(w.Body.Bytes(), actual); err != nil {
		t.Fatalf("failed to decode json: %v | %s", err, w.Body.Bytes())
	}

	testPayload(t, &repository.Payload{Error: "request is empty"}, actual)
}
