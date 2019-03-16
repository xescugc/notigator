package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"port"`

	Sources []Source `mapstructure:"sources"`
}

type Source struct {
	Name      string `mapstructure:"name"`
	Canonical string `mapstructure:"canonical"`

	Token  string `mapstructure:"token"`
	ApiKey string `mapstructure:"api_key"`
}

func New(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall the config file to the Config struct: %s", err)
	}

	return &cfg, nil
}
