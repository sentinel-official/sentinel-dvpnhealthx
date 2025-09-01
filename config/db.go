package config

import (
	"github.com/spf13/pflag"
)

type DBConfig struct {
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
}

// GetName returns the database name.
func (c *DBConfig) GetName() string {
	return c.Name
}

// GetURI returns the database URI.
func (c *DBConfig) GetURI() string {
	return c.URI
}

// SetForFlags registers command-line flags for DBConfig.
func (c *DBConfig) SetForFlags(f *pflag.FlagSet) {
	f.StringVar(&c.Name, "db.name", c.Name, "mongodb database name used for data storage")
	f.StringVar(&c.URI, "db.uri", c.URI, "mongodb connection uri including host and port")
}

func (c *DBConfig) Validate() error {
	return nil
}

// DefaultDBConfig returns a default database configuration instance.
func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Name: "dvpnhealthx",
		URI:  "mongodb://localhost:27017",
	}
}
