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
