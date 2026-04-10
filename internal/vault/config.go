package vault

import (
	"errors"
	"os"
)

// Config holds the configuration required to connect to Vault.
type Config struct {
	// Address is the base URL of the Vault server, e.g. "https://vault.example.com".
	Address string

	// Token is the Vault authentication token.
	Token string

	// SecretPath is the KV v2 path to read secrets from,
	// e.g. "secret/data/myapp/production".
	SecretPath string
}

// ConfigFromEnv builds a Config by reading standard Vault environment variables:
//   - VAULT_ADDR
//   - VAULT_TOKEN
//   - VAULT_SECRET_PATH
func ConfigFromEnv() (*Config, error) {
	cfg := &Config{
		Address:    os.Getenv("VAULT_ADDR"),
		Token:      os.Getenv("VAULT_TOKEN"),
		SecretPath: os.Getenv("VAULT_SECRET_PATH"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate returns an error if any required field is missing.
func (c *Config) Validate() error {
	if c.Address == "" {
		return errors.New("vault: VAULT_ADDR is required")
	}
	if c.Token == "" {
		return errors.New("vault: VAULT_TOKEN is required")
	}
	if c.SecretPath == "" {
		return errors.New("vault: VAULT_SECRET_PATH is required")
	}
	return nil
}

// NewClientFromConfig is a convenience constructor that creates a Client
// directly from a validated Config.
func NewClientFromConfig(cfg *Config) *Client {
	return NewClient(cfg.Address, cfg.Token)
}
