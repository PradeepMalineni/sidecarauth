// service.go
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	logger "sidecarauth/utility"
	"time"
)

// MakeRequest makes a request to the server
func MakeRequest(backendURL, authToken, httpMethod, contentType, payload string, headers http.Header) (string, error) {
	logger.Log("Service API : API Request Started ")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var req *http.Request
	var err error // Declare err variable
	//var req *http.Request
	//fmt.Println("BackendURL", backendURL)
	if httpMethod == "GET" {
		logger.Log("Service API : HTTP GET call")
		req, err = http.NewRequest(httpMethod, backendURL, nil)
		if err != nil {
			return "", fmt.Errorf("error creating request: %v", err)
		}
	} else {
		req, err = http.NewRequest(httpMethod, backendURL, bytes.NewBuffer([]byte(payload)))
		logger.Log("Service API : NOT HTTP GET call")
		if err != nil {
			return "", fmt.Errorf("error creating request: %v", err)
		}

		//req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Authorization", authToken)
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	// Make a request to the server using client.Do(req)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Format and print the JSON response
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON response: %v", err)
	}

	// Return the formatted response
	return prettyJSON.String(), nil
	// Handle the response
}
