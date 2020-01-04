package config

const (
	// DefaultPathToDatabase defines the path to the sqlite storage file
	DefaultPathToDatabase = "./storage"

	// DefaultPathToSchemaScripts defines the path to the scripts managing the database structure
	DefaultPathToSchemaScripts = "./scripts/schema"
)

// Database provides the configuration for the database
type Database struct {
	StoragePath       *string `json:"storage_path"`
	SchemaScriptsPath *string `json:"schema_scripts_path"`
}

// GetStoragePath returns the path to the storage file
func (d *Database) GetStoragePath() string {
	if d == nil || d.StoragePath == nil {
		return DefaultPathToDatabase
	}

	return *d.StoragePath
}

// GetSchemaScriptPath returns the path to the schema script files managing the database structure
func (d *Database) GetSchemaScriptPath() string {
	if d == nil || d.SchemaScriptsPath == nil {
		return DefaultPathToSchemaScripts
	}

	return *d.SchemaScriptsPath
}

// Merge overwrites the values given by the config from parameter if they differ from default values
func (d *Database) Merge(cfg *Database) {
	if cfg == nil || d == nil {
		return
	}

	if cfg.GetSchemaScriptPath() != DefaultPathToSchemaScripts {
		d.SchemaScriptsPath = cfg.SchemaScriptsPath
	}

	if cfg.GetStoragePath() != DefaultPathToDatabase {
		d.StoragePath = cfg.StoragePath
	}
}
