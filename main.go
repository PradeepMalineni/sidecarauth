// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sidecarauth/reverseproxy"
)

/*
func main() {
	// Initialize and use the auth sidecar proxy
	//authProxy := auth.NewAuthProxy()
	//authResult := authProxy.Authenticate("username", "password")
	fmt.Println("Authentication Result:", "hello2")
	fmt.Println("Authentication Result:", auth.X)

	// Initialize and use the cache module
	cacheInstance := cache.NewCache()
	cacheInstance.Set("key", "value")
	cachedValue, exists := cacheInstance.Get("key")
	if exists {
		fmt.Println("Cached Value:", cachedValue)
	} else {
		fmt.Println("Key not found in cache.")
	}
	fmt.Println("Starting the listening services for OAuth.")

	go_test.TestServices()

	//Start the Listner Function

	fmt.Println("Starting the Listner")
	listener.Listner()
}

*/

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

	//fmt.Println("Staring the test service on port 8443")
	//test_svcs.TestServices()

	// Create a reverse proxy that forwards requests to another service with a custom trust store
	proxy, err := reverseproxy.NewReverseProxy(cfg.TargetURL, cfg.TrustStoreFile)
	if err != nil {
		log.Fatal("Error creating reverse proxy:", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: proxy,
	}

	fmt.Printf("Proxy server listening on %s...\n", server.Addr)
	err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	//select {}
	auth.Initialize(config.AuthConfig.TokenURL, config.AuthConfig.AuthorizationHeader)
	http.HandleFunc("/listener-service", func(w http.ResponseWriter, r *http.Request) {
		tokenResponse, err := auth.GetAccessToken()
		if err != nil {
			http.Error(w, "Error getting access token", http.StatusInternalServerError)
			return
		}

		// Access fields from the tokenResponse as needed
		fmt.Println("Main-TokenType:", tokenResponse.TokenType)
		fmt.Println("Main-Access Token:", tokenResponse.AccessToken)
		fmt.Println("Main-Issued_at:", tokenResponse.IssuedAt)
		fmt.Println("Main-Expires In:", tokenResponse.ExpiresIn)
		fmt.Println("Main-Scope:", tokenResponse.Scope)

		responseJSON, err := json.Marshal(tokenResponse)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

	})

	// Specify the port to listen on
	port := 8090

	// Start the HTTP server
	fmt.Printf("Go HTTP Listener is listening on port %d...\n", port)
	//http.HandleFunc("/generate-token", handler)
	http.ListenAndServe(":8090", nil)

}
