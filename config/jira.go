package config

// Jira provides the configuration for Jira
type Jira struct {
	BaseURL  *string `json:"base_url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// GetBaseURL returns the base url
func (j *Jira) GetBaseURL() string {
	if j == nil {
		return ""
	}

	return *j.BaseURL
}

// GetUsername returns the username
func (j *Jira) GetUsername() string {
	if j == nil {
		return ""
	}

	return *j.Username
}

// GetPassword returns the username
func (j *Jira) GetPassword() string {
	if j == nil {
		return ""
	}

	return *j.Password
}
