// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sidecarauth/auth"
)

type AuthConfig struct {
	TokenURL            string `json:"TokenURL"`
	AuthorizationHeader string `json:"AuthorizationHeader"`
}

// Config struct for overall configuration
type Config struct {
	AuthConfig     AuthConfig  `json:"AuthConfig"`
	ListenerConfig interface{} `json:"ListenerConfig"` // Adjust the type as needed
	ServiceConfig  interface{} `json:"ServiceConfig"`  // Adjust the type as needed
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

func main() {
	configPath := "config.json"
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	auth.Initialize(config.AuthConfig.TokenURL, config.AuthConfig.AuthorizationHeader)
	tokenResponse, err := auth.GetAccessToken()
	if err != nil {
		log.Fatal(err) // Handle the error by logging and exiting
	}

	// Access fields from the tokenResponse as needed
	fmt.Println("Main-TokenType:", tokenResponse.TokenType)
	fmt.Println("Main-Access Token:", tokenResponse.AccessToken)
	fmt.Println("Main-Issued_at:", tokenResponse.IssuedAt)
	fmt.Println("Main-Expires In:", tokenResponse.ExpiresIn)
	fmt.Println("Main-Scope:", tokenResponse.Scope)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is the Go HTTP Listener!")
	})

	// Specify the port to listen on
	port := 8080

	// Start the HTTP server
	fmt.Printf("Go HTTP Listener is listening on port %d...\n", port)
	http.HandleFunc("/generate-token", handler)
	http.ListenAndServe(":8090", nil)

}
