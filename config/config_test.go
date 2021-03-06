package config_test

import (
	"errors"
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

	jiraBaseURL := "https://jira.atlassion.com"
	jiraUser := "jira"
	jiraPassword := "let me in"

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
				BaseURL:  &jiraBaseURL,
				Username: &jiraUser,
				Password: &jiraPassword,
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
		expected: &config.Config{},
		err:      config.ErrFileNotFound,
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcConfig{
		name:     "not a JSON format",
		filename: filepath.Join(".", "testdata", "test_config_error.json"),
		expected: &config.Config{},
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
			cfg := config.New()

			err := cfg.Load(testCase.filename)
			if !errors.Is(err, testCase.err) {
				t.Fatalf("unexpected error, expected '%v' but got '%v'", testCase.err, err)
			}

			testConfig(t, testCase.expected, cfg)
		})
	}
}

func testConfig(t *testing.T, expected, got *config.Config) {
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

	testDatabase(t, expected.GetDB(), got.GetDB())
	testGit(t, expected.GetGit(), got.GetGit())
	testJira(t, expected.GetJira(), got.GetJira())
	testService(t, expected.GetService(), got.GetService())
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

type tcConfigMerge struct {
	name      string
	actual    *config.Config
	mergeWith *config.Config
	expected  *config.Config
}

func getTestCasesConfigMerge(t *testing.T) []tcConfigMerge {
	t.Helper()

	var testCases []tcConfigMerge

	storagePath := "/mystorage"
	scriptPath := "/myscripts"
	db := &config.Database{
		StoragePath:       &storagePath,
		SchemaScriptsPath: &scriptPath,
	}

	baseURL := "another.url"
	branchPrefix := "myprefix"
	git := &config.Git{
		BaseURL:             &baseURL,
		ReleaseBranchPrefix: &branchPrefix,
	}

	username := "myUsername"
	password := "myPassword"
	jira := &config.Jira{
		BaseURL:  &baseURL,
		Username: &username,
		Password: &password,
	}

	port := 5000
	service := &config.Service{Port: &port}

	// 1.
	tc := tcConfigMerge{
		name:      "config nil",
		mergeWith: &config.Config{DB: db},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcConfigMerge{
		name:     "parameter nil",
		actual:   &config.Config{DB: db},
		expected: &config.Config{DB: db},
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcConfigMerge{
		name:      "config has default values, parameter has values",
		actual:    config.New(),
		mergeWith: &config.Config{DB: db, Git: git, Jira: jira, Service: service},
		expected:  &config.Config{DB: db, Git: git, Jira: jira, Service: service},
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestConfig_Merge(t *testing.T) {
	for _, testCase := range getTestCasesConfigMerge(t) {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.actual.Merge(testCase.mergeWith)
			testConfig(t, testCase.expected, testCase.actual)
		})
	}
}
