package config_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/rebel-l/branma_be/config"
)

type tcConfig struct {
	name     string
	filename string
	expected *config.Config
	err      error
}

func getTestCasesConfig() []tcConfig {
	var testCases []tcConfig

	// 1.
	storagePath := "./my_storage_path/"
	scriptPath := "./my_schema_script_path/"
	gitBaseURL := "https://github.com"
	gitPrefix := "live"
	port := 3333
	tc := tcConfig{
		name:     "success",
		filename: filepath.Join(".", "testdata", "test_config_success.json"),
		expected: &config.Config{
			DB: &config.Database{
				StoragePath:       &storagePath,
				SchemaScriptsPath: &scriptPath,
			},
			Git: &config.Git{
				BaseURL:             &gitBaseURL,
				ReleaseBranchPrefix: &gitPrefix,
			},
			Jira: &config.Jira{
				BaseURL:  "https://jira.atlassion.com",
				Username: "jira",
				Password: "let me in",
			},
			Service: &config.Service{
				Port: &port,
			},
		},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcConfig{
		name:     "no file",
		filename: filepath.Join(".", "testdata", "no_file.json"),
		err:      config.ErrFileNotFound,
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcConfig{
		name:     "not a JSON format",
		filename: filepath.Join(".", "testdata", "test_config_error.json"),
		err:      config.ErrNoJSONFormat,
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestConfig_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("long running test")
	}

	for _, testCase := range getTestCasesConfig() {
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

	if expected.GetDB().GetStoragePath() != got.GetDB().GetStoragePath() {
		t.Errorf("failed to set database storage path: expected '%s' but got '%s'",
			expected.GetDB().GetStoragePath(), got.GetDB().GetStoragePath())
	}

	if expected.GetDB().GetSchemaScriptPath() != got.GetDB().GetSchemaScriptPath() {
		t.Errorf("failed to set database schema script path: expected '%s' but got '%s'",
			expected.GetDB().GetSchemaScriptPath(), got.GetDB().GetSchemaScriptPath())
	}
}

func TestConfig_GetDB(t *testing.T) {
	var cfg *config.Config
	if cfg.GetDB() == nil {
		t.Errorf("failed to retrieve default value from nil struct")
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
