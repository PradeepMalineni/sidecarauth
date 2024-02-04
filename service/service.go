// service.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MakeRequest makes a request to the server
func MakeRequest(ctx context.Context, backendURL, authToken, httpMethod, contentType, payload string, headers http.Header) (string, error) {
	// Create a TLS configuration
	/*tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to true to skip server certificate verification
	}*/

	// Load the client certificate and private key
	/*cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return "", fmt.Errorf("error loading client certificate: %v", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if keyPassword != "" {
		// Add a callback for password retrieval if a password is specified
		tlsConfig.GetClientCertificate = func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}*/

	// Create a HTTP client with the custom TLS configuration
	/*client := &http.Client{
		Transport: &http.Transport{
			//TLSClientConfig: tlsConfig,
		},
	}
	*/
	client := &http.Client{}

	var req *http.Request
	var err error // Declare err variable
	//var req *http.Request
	//fmt.Println("BackendURL", backendURL)
	if httpMethod == "GET" {
		req, err = http.NewRequest(httpMethod, backendURL, nil)
		if err != nil {
			return "", fmt.Errorf("error creating request: %v", err)
		}
	} else {
		req, err = http.NewRequest(httpMethod, backendURL, bytes.NewBuffer([]byte(payload)))
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
