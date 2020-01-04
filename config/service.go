package config

const (
	// DefaultPort defines the default port the service listens on
	DefaultPort = 3000
)

// Service provides the configuration for service
type Service struct {
	Port *int `json:"port"`
}

// GetPort returns the port
func (s *Service) GetPort() int {
	if s == nil || s.Port == nil {
		return DefaultPort
	}

	return *s.Port
}

// Merge overwrites the values given by the config from parameter if they differ from default values
func (s *Service) Merge(cfg *Service) {
	if cfg == nil || s == nil {
		return
	}

	if cfg.GetPort() != DefaultPort {
		s.Port = cfg.Port
	}
}
