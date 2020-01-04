package config

// Jira provides the configuration for Jira
type Jira struct {
	BaseURL  *string `json:"base_url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// GetBaseURL returns the base url
func (j *Jira) GetBaseURL() string {
	if j == nil || j.BaseURL == nil {
		return ""
	}

	return *j.BaseURL
}

// GetUsername returns the username
func (j *Jira) GetUsername() string {
	if j == nil || j.Username == nil {
		return ""
	}

	return *j.Username
}

// GetPassword returns the username
func (j *Jira) GetPassword() string {
	if j == nil || j.Password == nil {
		return ""
	}

	return *j.Password
}

// Merge overwrites the values given by the config from parameter if they differ from default values
func (j *Jira) Merge(cfg *Jira) {
	if cfg == nil || j == nil {
		return
	}

	if cfg.GetBaseURL() != "" {
		j.BaseURL = cfg.BaseURL
	}

	if cfg.GetUsername() != "" {
		j.Username = cfg.Username
	}

	if cfg.GetPassword() != "" {
		j.Password = cfg.Password
	}
}
