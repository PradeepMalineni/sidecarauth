// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sidecarauth/auth"
	"sidecarauth/config"
	"sidecarauth/service"
)

func main() {
	configPath := "config.json"
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Create an instance of AuthHandler for each environment
	authHandlers := make(map[string]*auth.AuthHandler)
	for env, envConfig := range config.AuthConfig {
		authHandler := auth.NewAuthHandler(envConfig)
		authHandlers[env] = authHandler
	}

	http.HandleFunc(config.ListenerConfig.ListenerURI, func(w http.ResponseWriter, r *http.Request) {
		// Get the port from the request URL
		_, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			fmt.Printf("Error extracting port from host: %v\n", err)
			return
		}

		// Choose the environment based on the port
		var env string
		for e, p := range config.ListenerConfig.PortNumber {
			if p == port {
				env = e
				break
			}
		}

		// Check if environment is found
		if env == "" {
			fmt.Printf("No environment found for port: %s\n", port)
			return
		}

		// Initialize AuthHandler with configuration values
		authHandlers[env].Initialize()

		tokenResponse, err := authHandlers[env].GetAccessToken()
		if err != nil {
			http.Error(w, "Error getting access token", http.StatusInternalServerError)
			return
		}

		responseJSON, err := json.Marshal(tokenResponse)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		httpMethod := r.Method
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "NA"
		}
		uri := r.URL.Path
		backendURL := config.ServiceConfig[env].ApiURL + uri
		accessToken := "Bearer " + tokenResponse.AccessToken
		formattedResponse, err := service.MakeRequest(backendURL, accessToken, httpMethod, contentType, string(payload), r.Header)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		fmt.Fprintf(w, "\n\nFormatted Response: %s", formattedResponse)
	})

	// Start HTTP server for each configured port
	for env, port := range config.ListenerConfig.PortNumber {
		env, port := env, port // Capture variables for the goroutine
		go func() {
			fmt.Printf("Go HTTP Listener for %s is listening on port %s...\n", env, port)
			err := http.ListenAndServe(":"+port, nil)
			if err != nil {
				log.Fatalf("Error starting HTTP server for %s: %v", env, err)
			}
		}()
	}

	select {}
}
