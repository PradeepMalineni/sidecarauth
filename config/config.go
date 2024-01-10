// config.go
package config

import (
	"encoding/json"
	"fmt"
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
	AuthConfig     map[string]AuthConfig    `json:"AuthConfig"`
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

func CreateFileInDirectory(directory, filename string) error {
	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
		fmt.Printf("Created directory: %s\n", directory)
	}

	// Create or open the file within the specified directory
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("Created file: %s\n", filePath)
	return nil
}
