package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestDatabase_GetStoragePath(t *testing.T) {
	var db *config.Database
	if db.GetStoragePath() != config.DefaultPathToDatabase {
		t.Errorf(
			"failed to retrieve default value '%s' from nil struct but got '%s'",
			config.DefaultPathToDatabase,
			db.GetStoragePath(),
		)
	}
}

func TestDatabase_GetSchemaScriptPath(t *testing.T) {
	var db *config.Database
	if db.GetSchemaScriptPath() != config.DefaultPathToSchemaScripts {
		t.Errorf(
			"failed to retrieve default value '%s' from nil struct but got '%s'",
			config.DefaultPathToSchemaScripts,
			db.GetSchemaScriptPath(),
		)
	}
}

type tcDatabaseMerge struct {
	name      string
	actual    *config.Database
	mergeWith *config.Database
	expected  *config.Database
}

func getTestCasesDatabaseMerge(t *testing.T) []tcDatabaseMerge { // nolint:dupl,funlen
	t.Helper()

	var testCases []tcDatabaseMerge

	storagePath := "/mystorage"
	scriptPath := "/myscripts"

	newStoragePath := "/mynewstorage"
	newScriptPath := "/mynewscripts"

	// 1.
	tc := tcDatabaseMerge{
		name:      "config nil",
		mergeWith: &config.Database{StoragePath: &storagePath},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcDatabaseMerge{
		name:     "parameter nil",
		actual:   &config.Database{StoragePath: &storagePath},
		expected: &config.Database{StoragePath: &storagePath},
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcDatabaseMerge{
		name:      "config has default values, parameter has values",
		actual:    &config.Database{},
		mergeWith: &config.Database{StoragePath: &storagePath, SchemaScriptsPath: &scriptPath},
		expected:  &config.Database{StoragePath: &storagePath, SchemaScriptsPath: &scriptPath},
	}

	testCases = append(testCases, tc)

	// 4.
	tc = tcDatabaseMerge{
		name:      "config has values, parameter has values",
		actual:    &config.Database{StoragePath: &storagePath, SchemaScriptsPath: &scriptPath},
		mergeWith: &config.Database{StoragePath: &newStoragePath, SchemaScriptsPath: &newScriptPath},
		expected:  &config.Database{StoragePath: &newStoragePath, SchemaScriptsPath: &newScriptPath},
	}

	testCases = append(testCases, tc)

	// 5.
	tc = tcDatabaseMerge{
		name:      "config has values, parameter has default values",
		actual:    &config.Database{StoragePath: &storagePath, SchemaScriptsPath: &scriptPath},
		mergeWith: &config.Database{},
		expected:  &config.Database{StoragePath: &storagePath, SchemaScriptsPath: &scriptPath},
	}

	testCases = append(testCases, tc)

	// 6.
	tc = tcDatabaseMerge{
		name:      "config has default values, parameter has default values",
		actual:    &config.Database{},
		mergeWith: &config.Database{},
		expected:  &config.Database{},
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestDatabase_Merge(t *testing.T) {
	for _, testCase := range getTestCasesDatabaseMerge(t) {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.actual.Merge(testCase.mergeWith)
			testDatabase(t, testCase.expected, testCase.actual)
		})
	}
}

func testDatabase(t *testing.T, expected, got *config.Database) {
	t.Helper()

	if expected.GetStoragePath() != got.GetStoragePath() {
		t.Errorf("failed to set database storage path: expected '%s' but got '%s'",
			expected.GetStoragePath(), got.GetStoragePath())
	}

	if expected.GetSchemaScriptPath() != got.GetSchemaScriptPath() {
		t.Errorf("failed to set database schema script path: expected '%s' but got '%s'",
			expected.GetSchemaScriptPath(), got.GetSchemaScriptPath())
	}
}
