package config

// Jira provides the configuration for Jira
type Jira struct {
	BaseURL  string `json:"base_url"`
	Username string `json:"username"`
	Password string `json:"password"`
}
