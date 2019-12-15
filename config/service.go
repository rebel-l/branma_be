package config

// Service provides the configuration for service
type Service struct {
	Port        int    `json:"port"`
	StoragePath string `json:"storage_path"`
}

// GetPort returns the port
func (s *Service) GetPort() int {
	if s == nil {
		return 0
	}

	return s.Port
}

// GetStoragePath returns the path to the storage path
func (s *Service) GetStoragePath() string {
	if s == nil {
		return ""
	}

	return s.StoragePath
}
