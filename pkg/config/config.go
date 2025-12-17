package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	GRPCPort int    `yaml:"grpc_port"`
	HTTPPort int    `yaml:"http_port"`
	Host     string `yaml:"host"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			GRPCPort: 9090,
			HTTPPort: 8080,
			Host:     "0.0.0.0",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
