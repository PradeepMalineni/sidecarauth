// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	auth.Initialize(config.AuthConfig.TokenURL, config.AuthConfig.AuthorizationHeader)
	http.HandleFunc(config.ListenerConfig.ListenerURI, func(w http.ResponseWriter, r *http.Request) {
		tokenResponse, err := auth.GetAccessToken()
		if err != nil {
			http.Error(w, "Error getting access token", http.StatusInternalServerError)
			return
		}

		// Access fields from the tokenResponse as needed
		//fmt.Println("Main-TokenType:", tokenResponse.TokenType)
		//fmt.Println("Main-Access Token:", tokenResponse.AccessToken)
		//fmt.Println("Main-Issued_at:", tokenResponse.IssuedAt)
		//fmt.Println("Main-Expires In:", tokenResponse.ExpiresIn)
		//fmt.Println("Main-Scope:", tokenResponse.Scope)

		responseJSON, err := json.Marshal(tokenResponse)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header and write the JSON response
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
		//fmt.Println("IncomingURI", uri)
		formattedResponse, err := service.MakeRequest(config.ServiceConfig.ApiURL, uri, config.ServiceConfig.CertFile, config.ServiceConfig.KeyFile, tokenResponse.AccessToken, httpMethod, contentType, string(payload))
		if err != nil {
			// Handle the error as needed
			fmt.Println("Error making request:", err)
			return
		}
		fmt.Fprintf(w, "\n\nFormatted Response: %s", formattedResponse)
	})

	// Specify the port to listen on
	port := config.ListenerConfig.PortNumber

	// Start the HTTP server
	fmt.Printf("Go HTTP Listener is listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)

}
