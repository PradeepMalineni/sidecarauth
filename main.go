/*
Copyright (c) 2022-2024 All rights reserved.


For inquiries or collaboration, please contact:
- Pradeep Malineni <pradeep.malineni@hotmail.com>
- Bhumika Pehwani <bhumika15peshwani8@gmail.com>
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sidecarauth/auth"
	"sidecarauth/config"
	"sidecarauth/service"
	"time"
)

func main() {

	// Generate a timestamp
	currentTime := time.Now()
	// Format the timestamp as YYYY-MM-DD_HH-MM-SS
	timestampFormat := "2006-01-02"
	// Specify the directory name
	logDir := "logs"

	// Create the directory if it doesn't exist
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		fmt.Printf("Created directory: %s\n", logDir)
	}

	// Generate the log file name with the current timestamp
	logFileName := fmt.Sprintf("%s/app_%s.log", logDir, currentTime.Format(timestampFormat))

	// Open or create a log file
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[%s]: Error opening log file: %s", currentTime.Format(timestampFormat), err)
	}
	defer logFile.Close()

	// Set log output to both console and the log file
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Print("SideCarAuthSvcs Started")

	// Read the file path from the environment variable
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		fmt.Println("CONFIG_FILE_PATH environment variable not set.")
		os.Exit(1)
	}

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("[%s]: Error loading configuration %s:", currentTime.Format(timestampFormat), err)
	}
	// Create an instance of AuthHandler for each environment
	authHandlers := make(map[string]*auth.AuthHandler)
	//iterate over the config list
	for env, envConfig := range config.AuthConfig {
		authHandler := auth.NewAuthHandler(envConfig)
		authHandlers[env] = authHandler
	}
	log.Printf("[%s]: Authentication Listners enabled", currentTime.Format(timestampFormat))

	http.HandleFunc(config.ListenerConfig.ListenerURI, func(w http.ResponseWriter, r *http.Request) {
		// Get the port from the request URL
		_, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			fmt.Printf("[%s]: Error extracting port from host: %v\n", currentTime.Format(timestampFormat), err)
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
			fmt.Printf("[%s]: No environment found for port: %s\n", currentTime.Format(timestampFormat), port)
			return
		}

		// Initialize AuthHandler with configuration values
		authHandlers[env].Initialize()

		tokenResponse, err := newFunction(authHandlers, env)
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

func newFunction(authHandlers map[string]*auth.AuthHandler, env string) (auth.TokenResponse, error) {
	tokenResponse, err := authHandlers[env].GetAccessToken()
	return tokenResponse, err
}
