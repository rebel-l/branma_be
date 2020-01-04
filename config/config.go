package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	errFileNotFound = "configuration file not found"
	errLoadFile     = "failed to load configuration: %v"
	errNoJSONFormat = "content of file is not in JSON format"
	errReadData     = "failed to read data from file: %v"
)

var (
	// ErrFileNotFound is the error if the config file doesn't exist
	ErrFileNotFound = errors.New(errFileNotFound)

	// ErrNoJSONFormat indicates that the content is not a JSON
	ErrNoJSONFormat = errors.New(errNoJSONFormat)
)

// Config provides the complete configuration
type Config struct {
	DB      *Database `json:"db"`
	Git     *Git      `json:"git"`
	Jira    *Jira     `json:"jira"`
	Service *Service  `json:"service"`
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

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf(errReadData, err)
	}

	if err = json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("%w: %v", ErrNoJSONFormat, err)
	}

	return nil
}

//GetDB returns the configuration for the database
func (c *Config) GetDB() *Database {
	if c == nil {
		return &Database{}
	}

	return c.DB
}

// GetGit returns the configuration for Git
func (c *Config) GetGit() *Git {
	if c == nil {
		return &Git{}
	}

	return c.Git
}

// GetJira returns the configuration for Jira
func (c *Config) GetJira() *Jira {
	if c == nil {
		return &Jira{}
	}

	return c.Jira
}

// GetService returns the configuration for the service
func (c *Config) GetService() *Service {
	if c == nil {
		return &Service{}
	}

	return c.Service
}

// New returns the config with initialized sub configs
func New() *Config {
	return &Config{
		DB:      &Database{},
		Git:     &Git{},
		Jira:    &Jira{},
		Service: &Service{},
	}
}

// Merge overwrites the values given by the config from parameter if they differ from default values
func (c *Config) Merge(cfg *Config) {
	if cfg == nil || c == nil {
		return
	}

	c.GetDB().Merge(cfg.GetDB())
	c.GetGit().Merge(cfg.GetGit())
	c.GetJira().Merge(cfg.GetJira())
	c.GetService().Merge(cfg.GetService())
}
