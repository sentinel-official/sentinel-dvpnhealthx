package config

import (
	"fmt"

	"github.com/sentinel-official/sentinel-go-sdk/config"
	"github.com/spf13/pflag"
)

type Config struct {
	*config.Config `mapstructure:",squash"`
	API            *APIConfig `mapstructure:"api"`
	DB             *DBConfig  `mapstructure:"db"`
}

func (c *Config) SetForFlags(f *pflag.FlagSet) {
	c.Config.SetForFlags(f)
	c.API.SetForFlags(f)
	c.DB.SetForFlags(f)
}

func (c *Config) Validate() error {
	if err := c.Config.Validate(); err != nil {
		return fmt.Errorf("validating base config: %w", err)
	}
	if err := c.API.Validate(); err != nil {
		return fmt.Errorf("validating API config: %w", err)
	}
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("validating database config: %w", err)
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Config: config.DefaultConfig(),
		API:    DefaultAPIConfig(),
		DB:     DefaultDBConfig(),
	}
}
