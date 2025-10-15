package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// Kafka - contains all parameters kafka information.
type Kafka struct {
	Capacity uint64        `yaml:"capacity"`
	Topics   []string      `yaml:"topics"`
	GroupID  string        `yaml:"groupId"`
	Brokers  []string      `yaml:"brokers"`
	Tick     time.Duration `yaml:"tick"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project Project `yaml:"project"`
	Kafka   Kafka   `yaml:"kafka"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
