package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestService_GetPort(t *testing.T) {
	var service *config.Service
	if service.GetPort() != config.DefaultPort {
		t.Errorf(
			"failed to retrieve default value %d from nil struct, but got %d",
			config.DefaultPort,
			service.GetPort(),
		)
	}
}

func TestService_GetStoragePath(t *testing.T) {
	var service *config.Service
	if service.GetStoragePath() != config.DefaultPathToDatabase {
		t.Errorf(
			"failed to retrieve default value '%s' from nil struct but got '%s'",
			config.DefaultPathToDatabase,
			service.GetStoragePath(),
		)
	}
}

func TestService_GetSchemaScriptPath(t *testing.T) {
	var service *config.Service
	if service.GetSchemaScriptPath() != config.DefaultPathToSchemaScripts {
		t.Errorf(
			"failed to retrieve default value '%s' from nil struct but got '%s'",
			config.DefaultPathToSchemaScripts,
			service.GetSchemaScriptPath(),
		)
	}
}
