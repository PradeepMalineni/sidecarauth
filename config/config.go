package config

import (
	"encoding/json"
	"os"
)

// Configuration is a struct to hold the configuration parameters.
type Configuration struct {
	TargetURL      string `json:"targetURL"`
	ServerPort     int    `json:"serverPort"`
	CertFile       string `json:"certFile"`
	KeyFile        string `json:"keyFile"`
	TrustStoreFile string `json:"trustStoreFile"`
}

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
