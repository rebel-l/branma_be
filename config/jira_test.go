package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestJira_GetBaseURL(t *testing.T) {
	var jira *config.Jira
	if jira.GetBaseURL() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestJira_GetUsername(t *testing.T) {
	var jira *config.Jira
	if jira.GetUsername() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}

func TestJira_GetPassword(t *testing.T) {
	var jira *config.Jira
	if jira.GetPassword() != "" {
		t.Errorf("failed to retrieve default value from nil struct")
	}
}
