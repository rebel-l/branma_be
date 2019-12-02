package config

// Service provides the configuration for service
type Service struct {
	Port        int    `json:"port"`
	StoragePath string `json:"storage_path"`
}
