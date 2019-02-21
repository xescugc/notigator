package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"port"`
	GitHubToken  string `mapstructure:"github-token"`
	GitLabToken  string `mapstructure:"gitlab-token"`
	TrelloToken  string `mapstructure:"trello-token"`
	TrelloApiKey string `mapstructure:"trello-api-key"`
	TrelloMember string `mapstructure:"trello-api-key"`
}

func New(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall the config file to the Config struct: %s", err)
	}

	return &cfg, nil
}
