package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

const (
	defaultMaxUploadSizeMB = 2
)

type Config struct {
	// Spec config
	APIPath             string            `yaml:"api_path" envconfig:"API_PATH"`
	MediaPath           string            `yaml:"media_path" envconfig:"MEDIA_PATH"`
	AcceptedMimetypes   []string          `yaml:"accepted_mimetypes" envconfig:"ACCEPTED_MIMETYPES"`
	AllowAdultContent   bool              `yaml:"allow_adult_content" envconfig:"ALLOW_ADULT_CONTENT"`
	AllowViolentContent bool              `yaml:"allow_violent_content" envconfig:"ALLOW_VIOLENT_CONTENT"`
	Names               map[string]string `yaml:"names" envconfig:"NAMES"`

	// Server config
	Port            int               `yaml:"port" envconfig:"PORT"`
	StorageType     string            `yaml:"storage_type" envconfig:"STORAGE_TYPE"`
	StorageConfig   map[string]string `yaml:"storage_config" envconfig:"STORAGE_CONFIG"`
	MaxUploadSizeMB int64             `yaml:"max_upload_size_mb" envconfig:"MAX_UPLOAD_SIZE_MB"`
}

// Load Config from a yaml file at path.
func (c *Config) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	c.applyDefaults()

	return yaml.NewDecoder(f).Decode(c)
}

// Load Config from the environment.
func (c *Config) LoadFromEnv() error {
	c.applyDefaults()
	return envconfig.Process("", c)
}

func (c *Config) applyDefaults() {
	if c.MaxUploadSizeMB == 0 {
		c.MaxUploadSizeMB = defaultMaxUploadSizeMB
	}
}
