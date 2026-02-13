package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
	Logging  LoggingConfig  `json:"logging"`
	Cache    CacheConfig    `json:"cache"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	AuthPort     int    `json:"auth_port"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	UsersDir       string `json:"users_dir"`
	MaxConnections int    `json:"max_connections"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Cookie       string `json:"cookie"`
	Domain       string `json:"domain"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackUrl  string `json:"callback_url"`
	ReturnUrl    string `json:"return_url"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	DefaultExpiry   string `json:"default_expiry"`
	CleanupInterval string `json:"cleanup_interval"`
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Failed to load the env vars: %v", err)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the JSON
	jsonStr := os.ExpandEnv(string(data))

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.AuthPort == 0 {
		return fmt.Errorf("auth_port is required")
	}
	if c.Database.UsersDir == "" {
		return fmt.Errorf("users_dir is required")
	}
	if c.Auth.Cookie == "" {
		return fmt.Errorf("cookie authentication key is required")
	}
	if c.Auth.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if c.Auth.ClientID == "" {
		return fmt.Errorf("client id is required")
	}
	if c.Auth.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}
	if c.Auth.CallbackUrl == "" {
		return fmt.Errorf("callback url is required")
	}
	if c.Auth.ReturnUrl == "" {
		return fmt.Errorf("return url is required")
	}
	return nil
}

// ParseDuration parses a duration string
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
