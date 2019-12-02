package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	errFileNotFound = "configuration file not found"
	errLoadFile     = "failed to load configuration: %v"
	errNoJSONFormat = "content of file is not in JSON format: %v"
	errReadData     = "failed to read data from file: %v"
)

var (
	// ErrFileNotFound is the error if the config file doesn't exist
	ErrFileNotFound = errors.New(errFileNotFound)
)

// Config provides the complete configuration
type Config struct {
	Git     Git     `json:"git"`
	Jira    Jira    `json:"jira"`
	Service Service `json:"service"`
}

// Load loads the given JSON file into the struct
func (c *Config) Load(fileName string) error {
	fileName = filepath.Clean(fileName)
	if !osutils.FileOrPathExists(fileName) {
		return ErrFileNotFound
	}

	f, err := os.Open(fileName) // nolint: gosec
	if err != nil {
		return fmt.Errorf(errLoadFile, err)
	}

	defer func() {
		_ = f.Close()
	}()

	var data []byte
	if _, err = f.Read(data); err != nil {
		return fmt.Errorf(errReadData, err)
	}

	if err = json.Unmarshal(data, c); err != nil {
		return fmt.Errorf(errNoJSONFormat, err)
	}

	return nil
}

// New tries to load the config from file, merge it with the cli parameters
// and returns the final config.
func New(configFile string) (*Config, error) {
	cfg := &Config{}
	if err := cfg.Load(configFile); err != nil {
		return cfg, err
	}

	return cfg, nil
}
