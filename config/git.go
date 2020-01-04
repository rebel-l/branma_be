package config

// Git provides the configuration for Git
type Git struct {
	BaseURL             *string `json:"base_url"`
	ReleaseBranchPrefix *string `json:"release_branch_prefix"`
}

// GetBaseURL returns the base url
func (d *Git) GetBaseURL() string {
	if d == nil || d.BaseURL == nil {
		return ""
	}

	return *d.BaseURL
}

// GetReleaseBranchPrefix returns the prefix of the release branch
func (d *Git) GetReleaseBranchPrefix() string {
	if d == nil || d.ReleaseBranchPrefix == nil {
		return ""
	}

	return *d.ReleaseBranchPrefix
}

// Merge overwrites the values given by the config from parameter if they differ from default values
func (g *Git) Merge(cfg *Git) {
	if cfg == nil || g == nil {
		return
	}

	if cfg.GetReleaseBranchPrefix() != "" {
		g.ReleaseBranchPrefix = cfg.ReleaseBranchPrefix
	}

	if cfg.GetBaseURL() != "" {
		g.BaseURL = cfg.BaseURL
	}
}
