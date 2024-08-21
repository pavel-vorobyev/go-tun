package config

import (
	"encoding/json"
	"os"
)

type Provider interface {
	GetConfig() (*Config, error)
}

type DefaultConfigProvider struct{}

func (cp *DefaultConfigProvider) GetConfig() (*Config, error) {
	data, err := os.ReadFile("conf.json")
	if err != nil {
		return nil, err
	}

	var config *Config
	err = json.Unmarshal(data, &config)
	return config, err
}
