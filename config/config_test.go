package config_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestConfig_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	testCases := []struct {
		name     string
		filename string
		expected *config.Config
		err      error
	}{
		{
			name:     "success",
			filename: filepath.Join(".", "testdata", "test_config_success.json"),
			expected: &config.Config{
				Git: &config.Git{
					BaseURL:             "https://github.com",
					ReleaseBranchPrefix: "live",
				},
				Jira: &config.Jira{
					BaseURL:  "https://jira.atlassion.com",
					Username: "jira",
					Password: "let me in",
				},
				Service: &config.Service{
					Port:              3333,
					StoragePath:       "./my_storage_path/",
					SchemaScriptsPath: "./my_schema_script_path/",
				},
			},
		},

		{
			name:     "no file",
			filename: filepath.Join(".", "testdata", "no_file.json"),
			err:      config.ErrFileNotFound,
		},
		{
			name:     "not a JSON format",
			filename: filepath.Join(".", "testdata", "test_config_error.json"),
			err:      config.ErrNoJSONFormat,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := config.New(testCase.filename)
			fmt.Println(errors.Is(err, testCase.err))
			if !errors.Is(err, testCase.err) {
				t.Fatalf("unexpected error, expected '%v' but got '%v'", testCase.err, err)
			}

			assertConfig(t, testCase.expected, cfg)
		})
	}
}

func assertConfig(t *testing.T, expected, got *config.Config) {
	t.Helper()

	if expected == nil && got == nil {
		return
	}

	if expected == nil && got != nil {
		t.Fatalf("expected config to be nil but got '%#v'", got)
		return
	}

	if expected != nil && got == nil {
		t.Fatalf("expected config to be '%#v' but got nil", expected)
		return
	}

	// GIT
	if expected.GetGit().GetBaseURL() != got.GetGit().GetBaseURL() {
		t.Errorf("failed to set git base url: expected '%s' but got '%s'",
			expected.GetGit().GetBaseURL(), got.GetGit().GetBaseURL())
	}

	if expected.GetGit().GetReleaseBranchPrefix() != got.GetGit().GetReleaseBranchPrefix() {
		t.Errorf("failed to set git release branch prefix: expected '%s' but got '%s'",
			expected.GetGit().GetReleaseBranchPrefix(), got.GetGit().GetReleaseBranchPrefix())
	}

	// JIRA
	if expected.GetJira().GetBaseURL() != got.GetJira().GetBaseURL() {
		t.Errorf("failed to set JIRA base url: expected '%s' but got '%s'",
			expected.GetJira().GetBaseURL(), got.GetJira().GetBaseURL())
	}

	if expected.GetJira().GetUsername() != got.GetJira().GetUsername() {
		t.Errorf("failed to set JIRA username: expected '%s' but got '%s'",
			expected.GetJira().GetUsername(), got.GetJira().GetUsername())
	}

	if expected.GetJira().GetPassword() != got.GetJira().GetPassword() {
		t.Errorf("failed to set JIRA password: expected '%s' but got '%s'",
			expected.GetJira().GetPassword(), got.GetJira().GetPassword())
	}

	// Service
	if expected.GetService().GetPort() != got.GetService().GetPort() {
		t.Errorf("failed to set service port: expected '%d' but got '%d'",
			expected.GetService().GetPort(), got.GetService().GetPort())
	}

	if expected.GetService().GetStoragePath() != got.GetService().GetStoragePath() {
		t.Errorf("failed to set service storage path: expected '%s' but got '%s'",
			expected.GetService().GetStoragePath(), got.GetService().GetStoragePath())
	}

	if expected.GetService().GetSchemaScriptPath() != got.GetService().GetSchemaScriptPath() {
		t.Errorf("failed to set service schema script path: expected '%s' but got '%s'",
			expected.GetService().GetSchemaScriptPath(), got.GetService().GetSchemaScriptPath())
	}
}

func TestConfig_GetGit(t *testing.T) {
	var cfg *config.Config
	if cfg.GetGit() == nil {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestConfig_GetJira(t *testing.T) {
	var cfg *config.Config
	if cfg.GetJira() == nil {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestConfig_GetService(t *testing.T) {
	var cfg *config.Config
	if cfg.GetService() == nil {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}
