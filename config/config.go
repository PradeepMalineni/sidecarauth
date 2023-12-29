package config

import (
	"encoding/json"
	"os"
)

// Configuration is a struct to hold the configuration parameters.
/*type Configuration struct {
	TargetURL      string `json:"targetURL"`
	ServerPort     int    `json:"serverPort"`
	CertFile       string `json:"certFile"`
	KeyFile        string `json:"keyFile"`
	TrustStoreFile string `json:"trustStoreFile"`
}*/
type AuthConfig struct {
	TokenURL            string `json:"TokenURL"`
	AuthorizationHeader string `json:"AuthorizationHeader"`
}
type ListenerConfig struct {
	ListenerURI string `json:"ListenerURI"`
	PortNumber  string `json:"PortNumber"`
}
type ServiceConfig struct {
	ApiURL   string `json:"ApiURL"`
	CertFile string `json:"CertFile"`
	KeyFile  string `json:"KeyFile"`
}
type Config struct {
	AuthConfig     AuthConfig     `json:"AuthConfig"`
	ListenerConfig ListenerConfig `json:"ListenerConfig"` // Adjust the type as needed
	ServiceConfig  ServiceConfig  `json:"ServiceConfig"`  // Adjust the type as needed
}

/*
	func Load(filename string) (*Configuration, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		cfg := &Configuration{}
		err = decoder.Decode(cfg)
		if err != nil {
			return nil, err
		}

		return cfg, nil
	}
*/
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
