package config

// Git provides the configuration for Git
type Git struct {
	BaseURL             *string `json:"base_url"`
	ReleaseBranchPrefix *string `json:"release_branch_prefix"`
}

// GetBaseURL returns the base url
func (g *Git) GetBaseURL() string {
	if g == nil || g.BaseURL == nil {
		return ""
	}

	return *g.BaseURL
}

// GetReleaseBranchPrefix returns the prefix of the release branch
func (g *Git) GetReleaseBranchPrefix() string {
	if g == nil || g.ReleaseBranchPrefix == nil {
		return ""
	}

	return *g.ReleaseBranchPrefix
}
