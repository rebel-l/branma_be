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
	if s == nil {
		return DefaultPort
	}

	return *s.Port
}
