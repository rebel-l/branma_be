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
