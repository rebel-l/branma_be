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
