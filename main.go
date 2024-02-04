/*
Copyright (c) 2022-2024 All rights reserved.


For inquiries or collaboration, please contact:
- Pradeep Malineni <pradeep.malineni@hotmail.com>
- Bhumika Pehwani <bhumika15peshwani8@gmail.com>
*/
//testing for collaboration
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sidecarauth/auth"
	"sidecarauth/config"
	"sidecarauth/service"
	logger "sidecarauth/utility"
	"time"
)

func main() {

	// Generate a timestamp

	// Format the timestamp as YYYY-MM-DD_HH-MM-SS

	// Specify the directory name
	//updated the code here only
	logger.InitLogger()

	// Read the file path from the environment variable
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		logger.Log("Error : CONFIG_FILE_PATH environment variable not set. Program Exit")
		os.Exit(1)
	}
	logger.Log("Sidecard Authentication service started ....")

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.LogF("Error : loading configuration", err)
	}
	logger.LogF("Configuration file loaded sucessfully", configFilePath)
	logger.Log("Authentication handler initiation started")

	// Create an instance of AuthHandler for each environment
	authHandlers := make(map[string]*auth.AuthHandler)
	//iterate over the config list
	for env, envConfig := range config.AuthConfig {
		authHandler := auth.NewAuthHandler(env, envConfig)
		authHandlers[env] = authHandler
	}

	http.HandleFunc(config.ListenerConfig.ListenerURI, func(w http.ResponseWriter, r *http.Request) {
		// Get the port from the request URL
		logger.Log("HTTP Handler function execution started")

		_, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			logger.LogF("Error extracting port from host: ", err)
			os.Exit(1)
		}

		// Choose the environment based on the port
		var env string
		for e, p := range config.ListenerConfig.PortNumber {
			if p == port {
				env = e
				logger.LogF("Environment  ", env)
				break
			}
		}

		// Check if environment is found
		if env == "" {
			logger.LogF("No environment found for port:", port)
			return
		}
		logger.LogF("authHandlers Initalizing the token request for  ", env)

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

		//Add logic not execute the followimng code without token response

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
		logger.Log("Initating the API Request")

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
			server := &http.Server{
				Addr: ":" + port,
				// Add other server configuration options as needed
			}

			// Use a timeout for reading requests from the client
			server.ReadTimeout = 10 * time.Second // Adjust the timeout duration as needed

			// Use a timeout for writing responses to the client
			server.WriteTimeout = 10 * time.Second // Adjust the timeout duration as needed

			server.IdleTimeout = 10 * time.Second

			err := server.ListenAndServe()
			if err != nil {
				logger.LogF("Error starting HTTP server for", env)
			}
		}()
	}

	select {}
}

func newFunction(authHandlers map[string]*auth.AuthHandler, env string) (auth.TokenResponse, error) {
	logger.Log("newFuncion check the token respose")
	tokenResponse, err := authHandlers[env].GetAccessToken()
	return tokenResponse, err
}
