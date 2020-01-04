package config

const (
	// DefaultPort defines the default port the service listens on
	DefaultPort = 3000

	// DefaultPathToDatabase defines the path to the sqlite storage file
	DefaultPathToDatabase = "./storage"

	// DefaultPathToSchemaScripts defines the path to the scripts managing the database structure
	DefaultPathToSchemaScripts = "./scripts/schema"
)

// Service provides the configuration for service
type Service struct {
	Port              int    `json:"port"`
	StoragePath       string `json:"storage_path"`
	SchemaScriptsPath string `json:"schema_scripts_path"`
}

// GetPort returns the port
func (s *Service) GetPort() int {
	if s == nil {
		return DefaultPort
	}

	return s.Port
}

// GetStoragePath returns the path to the storage file
func (s *Service) GetStoragePath() string {
	if s == nil {
		return DefaultPathToDatabase
	}

	return s.StoragePath
}

// GetSchemaScriptPath returns the path to the schema script files managing the database structure
func (s *Service) GetSchemaScriptPath() string {
	if s == nil {
		return DefaultPathToSchemaScripts
	}

	return s.SchemaScriptsPath
}
