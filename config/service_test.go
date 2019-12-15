package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestService_GetPort(t *testing.T) {
	var service *config.Service
	if service.GetPort() != 0 {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestService_GetStoragePath(t *testing.T) {
	var service *config.Service
	if service.GetStoragePath() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}
