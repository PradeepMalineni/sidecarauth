// auth.go
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	tokenResponse TokenResponse
	configURL     string
	authHeader    string
	mu            sync.Mutex // Mutex for thread-safe operations
)

// Initialize is called with the configuration values
func Initialize(url, header string) {
	configURL = url
	authHeader = header

	// Call the function to get the initial access token
	getAccessToken()
}

// GetAccessToken function to get the access token
func GetAccessToken() (TokenResponse, error) {
	// Lock to ensure thread-safe access
	mu.Lock()
	defer mu.Unlock()

	// Check if the token is expired or about to expire

	now := time.Now().Unix()
	//fmt.Println("now time", now)
	//fmt.Println("ExpiryTime", tokenResponse.ExpiresIn)
	//fmt.Println("difference left", tokenResponse.ExpiresIn-now)
	if now >= tokenResponse.ExpiresIn {
		// Token is expired or about to expire, refresh it
		getAccessToken()
	}

	return tokenResponse, nil
}

func getAccessToken() {
	// Use configURL and authHeader as needed
	//url := "https://apiidp-enterprise1-sandbox.wellsfargo.com/oauth/token"

	/*payload := strings.NewReader("grant_type=client_credentials")

	req, _ := http.NewRequest("POST", configURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic 123aqbc")

	res, _ := http.DefaultClient.Do(req)*/
	// Below line added for testing purpose, uncomment line 58 to 65 and comment 67 when ready.
	res, err := http.Get(configURL)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if err != nil {
		fmt.Println("Error performing HTTP request:", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK status:", res.Status)
		handleError(body)
		return
	}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	//fmt.Println("Refreshed access token:", tokenResponse.AccessToken)
}

func handleError(body []byte) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err == nil && errorResponse.Error != "" {
		fmt.Println("Error response from the server:", errorResponse.Error)
	} else {
		fmt.Println("Unexpected response from the server")
	}
}
