package config

import (
	"github.com/spf13/pflag"
)

type APIConfig struct {
	Addr string `mapstructure:"addr"`
}

func (c *APIConfig) GetAddr() string {
	return c.Addr
}

// SetForFlags registers command-line flags for APIConfig.
func (c *APIConfig) SetForFlags(f *pflag.FlagSet) {
	f.StringVar(&c.Addr, "api.addr", c.Addr, "listen address for api server communication")
}

func (c *APIConfig) Validate() error {
	return nil
}

// DefaultAPIConfig returns a default configuration instance.
func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		Addr: ":8080",
	}
}
