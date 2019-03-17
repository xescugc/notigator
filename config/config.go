package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the entity that holds all the
// configuration of the config file
type Config struct {
	Port string `mapstructure:"port"`

	Sources []Source `mapstructure:"sources"`
}

// Source it's the configuration of each one of
// the supported sources
type Source struct {
	Name      string `mapstructure:"name"`
	Canonical string `mapstructure:"canonical"`

	Token  string `mapstructure:"token"`
	APIKey string `mapstructure:"api_key"`
}

// New initializes a Config from the v
func New(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall the config file to the Config struct: %s", err)
	}

	return &cfg, nil
}
