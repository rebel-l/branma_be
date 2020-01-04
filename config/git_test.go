package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestGit_GetBaseURL(t *testing.T) {
	var git *config.Git
	if git.GetBaseURL() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestGit_GetReleaseBranchPrefix(t *testing.T) {
	var git *config.Git
	if git.GetReleaseBranchPrefix() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

type tcGitMerge struct {
	name      string
	actual    *config.Git
	mergeWith *config.Git
	expected  *config.Git
}

func getTestCasesGitMerge(t *testing.T) []tcGitMerge { // nolint:dupl,funlen
	t.Helper()

	var testCases []tcGitMerge

	baseURL := "git.url"
	branchPrefix := "myprefix"

	newBaseURL := "mynew.url"
	newBranchPrefix := "mynewprefix"

	// 1.
	tc := tcGitMerge{
		name:      "config nil",
		mergeWith: &config.Git{BaseURL: &baseURL},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcGitMerge{
		name:     "parameter nil",
		actual:   &config.Git{BaseURL: &baseURL},
		expected: &config.Git{BaseURL: &baseURL},
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcGitMerge{
		name:      "config has default values, parameter has values",
		actual:    &config.Git{},
		mergeWith: &config.Git{BaseURL: &baseURL, ReleaseBranchPrefix: &branchPrefix},
		expected:  &config.Git{BaseURL: &baseURL, ReleaseBranchPrefix: &branchPrefix},
	}

	testCases = append(testCases, tc)

	// 4.
	tc = tcGitMerge{
		name:      "config has values, parameter has values",
		actual:    &config.Git{BaseURL: &baseURL, ReleaseBranchPrefix: &branchPrefix},
		mergeWith: &config.Git{BaseURL: &newBaseURL, ReleaseBranchPrefix: &newBranchPrefix},
		expected:  &config.Git{BaseURL: &newBaseURL, ReleaseBranchPrefix: &newBranchPrefix},
	}

	testCases = append(testCases, tc)

	// 5.
	tc = tcGitMerge{
		name:      "config has values, parameter has default values",
		actual:    &config.Git{BaseURL: &baseURL, ReleaseBranchPrefix: &branchPrefix},
		mergeWith: &config.Git{},
		expected:  &config.Git{BaseURL: &baseURL, ReleaseBranchPrefix: &branchPrefix},
	}

	testCases = append(testCases, tc)

	// 6.
	tc = tcGitMerge{
		name:      "config has default values, parameter has default values",
		actual:    &config.Git{},
		mergeWith: &config.Git{},
		expected:  &config.Git{},
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestGit_Merge(t *testing.T) {
	for _, testCase := range getTestCasesGitMerge(t) {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.actual.Merge(testCase.mergeWith)
			testGit(t, testCase.expected, testCase.actual)
		})
	}
}

func testGit(t *testing.T, expected, got *config.Git) {
	t.Helper()

	if expected.GetBaseURL() != got.GetBaseURL() {
		t.Errorf("failed to set git base url: expected '%s' but got '%s'",
			expected.GetBaseURL(), got.GetBaseURL())
	}

	if expected.GetReleaseBranchPrefix() != got.GetReleaseBranchPrefix() {
		t.Errorf("failed to set git release branch prefix: expected '%s' but got '%s'",
			expected.GetReleaseBranchPrefix(), got.GetReleaseBranchPrefix())
	}
}
