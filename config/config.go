// config.go
package config

import (
	"encoding/json"
	"os"
)

type AuthConfig struct {
	TokenURL            string `json:"TokenURL"`
	AuthorizationHeader string `json:"AuthorizationHeader"`
}

type ListenerConfig struct {
	ListenerURI string            `json:"ListenerURI"`
	PortNumber  map[string]string `json:"PortNumber"` // Change type to map[string]string
}

type ServiceConfig struct {
	ApiURL string `json:"ApiURL"`
}

type Config struct {
	AuthConfig     AuthConfig               `json:"AuthConfig"`
	ListenerConfig ListenerConfig           `json:"ListenerConfig"`
	ServiceConfig  map[string]ServiceConfig `json:"ServiceConfig"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
