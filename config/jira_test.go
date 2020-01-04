package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestJira_GetBaseURL(t *testing.T) {
	var jira *config.Jira
	if jira.GetBaseURL() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestJira_GetUsername(t *testing.T) {
	var jira *config.Jira
	if jira.GetUsername() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestJira_GetPassword(t *testing.T) {
	var jira *config.Jira
	if jira.GetPassword() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

type tcJiraMerge struct {
	name      string
	actual    *config.Jira
	mergeWith *config.Jira
	expected  *config.Jira
}

func getTestCasesJiraMerge(t *testing.T) []tcJiraMerge { // nolint:funlen
	t.Helper()

	var testCases []tcJiraMerge

	baseURL := "my.url"
	username := "myUsername"
	password := "myPassword"

	newBaseURL := "my.url"
	newUsername := "myNewUsername"
	newPassword := "myNewPassword"

	// 1.
	tc := tcJiraMerge{
		name:      "config nil",
		mergeWith: &config.Jira{BaseURL: &baseURL},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcJiraMerge{
		name:     "parameter nil",
		actual:   &config.Jira{BaseURL: &baseURL},
		expected: &config.Jira{BaseURL: &baseURL},
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcJiraMerge{
		name:      "config has default values, parameter has values",
		actual:    &config.Jira{},
		mergeWith: &config.Jira{BaseURL: &baseURL, Username: &username, Password: &password},
		expected:  &config.Jira{BaseURL: &baseURL, Username: &username, Password: &password},
	}

	testCases = append(testCases, tc)

	// 4.
	tc = tcJiraMerge{
		name:      "config has values, parameter has values",
		actual:    &config.Jira{BaseURL: &baseURL, Username: &username, Password: &password},
		mergeWith: &config.Jira{BaseURL: &newBaseURL, Username: &newUsername, Password: &newPassword},
		expected:  &config.Jira{BaseURL: &newBaseURL, Username: &newUsername, Password: &newPassword},
	}

	testCases = append(testCases, tc)

	// 5.
	tc = tcJiraMerge{
		name:      "config has values, parameter has default values",
		actual:    &config.Jira{BaseURL: &baseURL, Username: &username, Password: &password},
		mergeWith: &config.Jira{},
		expected:  &config.Jira{BaseURL: &baseURL, Username: &username, Password: &password},
	}

	testCases = append(testCases, tc)

	// 6.
	tc = tcJiraMerge{
		name:      "config has default values, parameter has default values",
		actual:    &config.Jira{},
		mergeWith: &config.Jira{},
		expected:  &config.Jira{},
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestJira_Merge(t *testing.T) {
	for _, testCase := range getTestCasesJiraMerge(t) {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.actual.Merge(testCase.mergeWith)
			testJira(t, testCase.expected, testCase.actual)
		})
	}
}

func testJira(t *testing.T, expected, got *config.Jira) {
	t.Helper()

	if expected.GetBaseURL() != got.GetBaseURL() {
		t.Errorf("failed to set JIRA base url: expected '%s' but got '%s'",
			expected.GetBaseURL(), got.GetBaseURL())
	}

	if expected.GetUsername() != got.GetUsername() {
		t.Errorf("failed to set JIRA username: expected '%s' but got '%s'",
			expected.GetUsername(), got.GetUsername())
	}

	if expected.GetPassword() != got.GetPassword() {
		t.Errorf("failed to set JIRA password: expected '%s' but got '%s'",
			expected.GetPassword(), got.GetPassword())
	}
}
